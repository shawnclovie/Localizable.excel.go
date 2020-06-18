package excel

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/tealeg/xlsx"
)

const (
	exportFormatJSON       = "json"
	exportFormatIOSStrings = "ios_strings"
	exportFormatAndroidXML = "android_xml"
)

func IsExportFormatValid(f string) bool {
	switch strings.ToLower(f) {
	case exportFormatJSON, exportFormatIOSStrings:
		return true
	}
	return false
}

type Documents struct {
	Path      string
	Documents []*Document
}

func (docs *Documents) Export(basepath string) (err error) {
	for i, doc := range docs.Documents {
		switch doc.Format {
		case exportFormatJSON:
			err = doc.exportAsJSON(basepath)
		case exportFormatIOSStrings:
			err = doc.exportAsIOSStrings(basepath)
		case exportFormatAndroidXML:
			err = fmt.Errorf("not support %v yet", doc.Format)
		default:
			err = fmt.Errorf("invalid format %v", doc.Format)
		}
		if err != nil {
			return fmt.Errorf("export document(%v) failure: %w", i, err)
		}
	}
	return nil
}

func LoadExcelDocumentsFromFile(path string) (*Documents, error) {
	xls, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}
	docs := &Documents{Path: path}
	for sheetIndex, sheet := range xls.Sheets {
		if sheet.MaxCol <= 1 || sheet.MaxRow <= 1 {
			continue
		}
		meta, err := url.ParseQuery(strings.ReplaceAll(sheet.Cell(0, 0).Value, "\n", "&"))
		if err != nil {
			return nil, fmt.Errorf("sheet(%v) meta cell (0,0) parse failed %w", sheetIndex, err)
		}
		row0 := sheet.Rows[0]
		cellCount := len(row0.Cells)
		if cellCount != sheet.MaxCol {
			fmt.Printf("sheet(%v) count of row0.Cells(%v) != %v\n", sheetIndex, cellCount, sheet.MaxCol)
		}
		doc := &Document{Name: sheet.Name, Path: meta.Get("path"), Format: meta.Get("format")}
		if doc.Path == "" {
			return nil, fmt.Errorf("sheet(%v) no path defined in meta", sheetIndex)
		}
		if !IsExportFormatValid(doc.Format) {
			return nil, fmt.Errorf("sheet(%v) no valid format defined in meta (current is '%s')", sheetIndex, doc.Format)
		}
		doc.LanguageNames = make([]string, len(row0.Cells)-1)
		for index, cell := range row0.Cells {
			if index > 0 {
				doc.LanguageNames[index-1] = cell.Value
			}
		}
		keyCount := sheet.MaxRow - 1
		keys := make([]string, keyCount)
		for i := 0; i < sheet.MaxRow; i += 1 {
			if i > 0 {
				keys[i-1] = sheet.Cell(i, 0).Value
			}
		}
		doc.SetKeys(keys)

		for ci := 1; ci < sheet.MaxCol; ci += 1 {
			lang := doc.LanguageNames[ci-1]
			values := make([]string, keyCount)
			for ri := 1; ri < sheet.MaxRow; ri += 1 {
				cell := sheet.Cell(ri, ci)
				values[ri-1] = cell.Value
			}
			doc.set(lang, values)
		}
		docs.Documents = append(docs.Documents, doc)
	}
	return docs, nil
}

type Document struct {
	Name          string
	Path          string
	Format        string
	LanguageNames []string
	KeyNames      []string
	keysMap       map[uint32]int
	translations  map[string][]string
}

func (d *Document) SetKeys(keys []string) {
	d.KeyNames = keys
	d.keysMap = make(map[uint32]int, len(keys))
	for i, key := range keys {
		d.keysMap[hash(key)] = i
	}
}

func (d *Document) prepareTranslations() {
	if d.translations == nil {
		d.translations = make(map[string][]string, len(d.LanguageNames))
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (d *Document) StringsForLanguage(lang string) []string {
	return d.translations[lang]
}

func (d *Document) StringForLanguageAndKey(lang, key string) (string, bool) {
	keyIndex, found := d.keysMap[hash(key)]
	if !found {
		return "", false
	}
	values, found := d.translations[lang]
	if !found || len(values) < keyIndex {
		return "", false
	}
	return values[keyIndex], true
}

func (d *Document) set(language string, values []string) {
	d.prepareTranslations()
	d.translations[language] = values
}

func (d *Document) pathComponents() (dir string, file string) {
	slashPos := strings.LastIndex(d.Path, "/")
	if slashPos <= 0 {
		dir = "."
	} else {
		dir = d.Path[:slashPos]
	}
	if slashPos >= 0 {
		file = d.Path[slashPos+1:]
	} else {
		file = d.Path
	}
	return dir, file
}

func (doc *Document) exportAsIOSStrings(basepath string) error {
	for _, lang := range doc.LanguageNames {
		dir, file := doc.pathComponents()
		stringsPath := basepath + "/" + dir + "/" + lang + ".lproj/" + file + ".strings"
		translations := doc.stringMapForLanguage(lang)
		if len(translations) == 0 {
			println("translations for", lang, "is empty.")
			continue
		}
		println("writing document", doc.Name, "(", len(doc.KeyNames), "keys) to\n", stringsPath)
		dir = path.Dir(stringsPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0700)
		}
		var err error
		content := strings.Builder{}
		for _, key := range doc.KeyNames {
			text := translations[key]
			if strings.Contains(key, " ") {
				key = "\"" + key + "\""
			}
			text = strings.Replace(text, "\"", "\\\"", -1)
			text = strings.Replace(text, "\n", "\\n", -1)
			line := key + "=\"" + text + "\";\n"
			_, err = content.WriteString(line)
			if err != nil {
				break
			}
		}
		err = ioutil.WriteFile(stringsPath, []byte(content.String()), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func (doc *Document) stringMapForLanguage(lang string) map[string]string {
	m := make(map[string]string, len(doc.KeyNames))
	for _, key := range doc.KeyNames {
		v, found := doc.StringForLanguageAndKey(lang, key)
		if found {
			m[key] = v
		}
	}
	return m
}

func (doc *Document) exportAsJSON(basepath string) error {
	dir, file := doc.pathComponents()
	dir = basepath + "/" + dir
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}
	translations := make(map[string][]string, len(doc.LanguageNames))
	for _, lang := range doc.LanguageNames {
		m := doc.StringsForLanguage(lang)
		if len(m) > 0 {
			translations[lang] = m
		}
	}
	bs, err := json.Marshal(map[string]interface{}{
		"keys": doc.KeyNames,
		"text": translations,
	})
	if err != nil {
		return err
	}
	exportPath := dir + "/" + file + ".json"
	println("writing document", doc.Name, "(", len(doc.KeyNames), "keys) to\n", exportPath)
	return ioutil.WriteFile(exportPath, bs, 0666)
}
