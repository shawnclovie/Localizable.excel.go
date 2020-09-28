package excel

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/shawnclovie/Localizable.excel.go/utility"
	"github.com/tealeg/xlsx"
	"gopkg.in/yaml.v3"
)

const (
	exportFormatJSON    = "json"
	exportFormatIOS     = "ios"
	exportFormatAndroid = "android"
	exportFormatARB     = "arb"

	exportFileMode os.FileMode = 0644
)

type Documents struct {
	Path      string
	Documents []Document
}

func NewDocument(languages, docNames []string) (docs Documents) {
	for _, name := range docNames {
		d := Document{
			Name:          name,
			LanguageNames: languages,
			KeyNames:      []string{"first_key"},
		}
		for _, lang := range languages {
			d.set(lang, []string{""})
		}
		docs.Documents = append(docs.Documents, d)
	}
	return
}

func (docs *Documents) Export(basepath string) (err error) {
	for i, doc := range docs.Documents {
		doc.prepareKeysMap()
		switch doc.Format {
		case exportFormatJSON:
			err = doc.exportAsJSON(basepath)
		case exportFormatIOS:
			err = doc.exportAsIOSStrings(basepath)
		case exportFormatARB:
			err = doc.exportAsARB(basepath)
		case exportFormatAndroid:
			err = doc.exportAsAndroidStrings(basepath)
		default:
			err = fmt.Errorf("invalid format %v", doc.Format)
		}
		if err != nil {
			return fmt.Errorf("export document(%v) failure: %w", i, err)
		}
	}
	return nil
}

func (docs Documents) ToJSONData() (bs []byte, err error) {
	result := make([]interface{}, len(docs.Documents))
	for j, doc := range docs.Documents {
		result[j] = doc.toMap()
	}
	return json.MarshalIndent(result, "", "  ")
}

func (docs Documents) ToYAMLData() (bs []byte, err error) {
	result := make([]interface{}, len(docs.Documents))
	for j, doc := range docs.Documents {
		result[j] = doc.toMap()
	}
	return yaml.Marshal(result)
}

func (docs Documents) ToExcelData() (bs []byte, err error) {
	xls := xlsx.NewFile()
	tdStyle := xlsx.NewStyle()
	tdStyle.Font.Name = "Consolas"
	tdStyle.Font.Size = 10
	thStyle := xlsx.NewStyle()
	thStyle.Font.Bold = true
	thStyle.Font.Name = tdStyle.Font.Name
	thStyle.Font.Size = tdStyle.Font.Size
	for _, doc := range docs.Documents {
		sheet, sheetErr := xls.AddSheet(doc.Name)
		if sheetErr != nil {
			err = sheetErr
			return
		}
		setExcelCell(sheet.Cell(0, 0), fmt.Sprintf("path=%s&format=%s", doc.Path, doc.Format), tdStyle)
		for i, lang := range doc.LanguageNames {
			setExcelCell(sheet.Cell(0, i+1), lang, thStyle)
		}
		for i, key := range doc.KeyNames {
			setExcelCell(sheet.Cell(i+1, 0), key, thStyle)
		}
		for i, lang := range doc.LanguageNames {
			trans := doc.Translations[lang]
			for j, tran := range trans {
				setExcelCell(sheet.Cell(j+1, i+1), tran, tdStyle)
			}
		}
		sheet.SetColWidth(0, len(doc.LanguageNames), 24)
	}
	var buf bytes.Buffer
	err = xls.Write(&buf)
	bs = buf.Bytes()
	return
}

func setExcelCell(cell *xlsx.Cell, content string, style *xlsx.Style) {
	cell.SetString(content)
	cell.SetStyle(style)
}

func LoadDocumentsFromFile(path string) (docs Documents, err error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	var maps []map[string]interface{}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml":
		err = yaml.Unmarshal(bs, &maps)
	case ".json":
		err = json.Unmarshal(bs, &maps)
	default:
		err = fmt.Errorf("unknown extension: %v", ext)
	}
	if err != nil {
		return
	}
	for _, m := range maps {
		d := parseDocumentFromMap(m)
		docs.Documents = append(docs.Documents, d)
	}
	return
}

func LoadDocumentsFromExcelFile(path string) (docs Documents, err error) {
	xls, err := xlsx.OpenFile(path)
	if err != nil {
		return
	}
	docs.Path = path
	for sheetIndex, sheet := range xls.Sheets {
		if sheet.MaxCol <= 1 || sheet.MaxRow <= 1 {
			continue
		}
		meta, queryErr := url.ParseQuery(strings.ReplaceAll(sheet.Cell(0, 0).Value, "\n", "&"))
		if queryErr != nil {
			err = fmt.Errorf("sheet(%v) meta cell (0,0) parse failed %w", sheetIndex, queryErr)
			return
		}
		row0 := sheet.Rows[0]
		cellCount := len(row0.Cells)
		if cellCount != sheet.MaxCol {
			fmt.Printf("sheet(%v) count of row0.Cells(%v) != %v\n", sheetIndex, cellCount, sheet.MaxCol)
		}
		doc := Document{Name: sheet.Name, Path: meta.Get("path"), Format: meta.Get("format")}
		if doc.Path == "" {
			err = fmt.Errorf("sheet(%v) no path defined in meta", sheetIndex)
			return
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
		doc.KeyNames = keys

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
	return
}

type Document struct {
	Name          string
	Path          string
	Format        string
	LanguageNames []string
	KeyNames      []string
	Translations  map[string][]string

	keysMap map[uint32]int
}

func parseDocumentFromMap(encoded map[string]interface{}) (d Document) {
	d.Name = utility.AnyToString(encoded["name"])
	d.Path = utility.AnyToString(encoded["path"])
	d.Format = utility.AnyToString(encoded["format"])
	d.LanguageNames = utility.AnyToStringArray(encoded["language_names"])
	d.Translations = make(map[string][]string)
	lines := encoded["translations"].([]interface{})
	d.KeyNames = make([]string, len(lines))
	for j, it := range lines {
		line := utility.AnyToStringArray(it)
		d.KeyNames[j] = line[0]
		for i, lang := range d.LanguageNames {
			if d.Translations[lang] == nil {
				d.Translations[lang] = make([]string, len(lines))
			}
			d.Translations[lang][j] = line[i+1]
		}
	}
	return
}

func (d Document) toMap() map[string]interface{} {
	trans := make([][]string, len(d.KeyNames))
	for j, key := range d.KeyNames {
		line := make([]string, 1, len(d.LanguageNames)+1)
		line[0] = key
		for _, lang := range d.LanguageNames {
			line = append(line, d.Translations[lang][j])
		}
		trans[j] = line
	}
	return map[string]interface{}{
		"name":           d.Name,
		"path":           d.Path,
		"format":         d.Format,
		"language_names": d.LanguageNames,
		"translations":   trans,
	}
}

func (d *Document) prepareKeysMap() {
	d.keysMap = make(map[uint32]int, len(d.KeyNames))
	for i, key := range d.KeyNames {
		d.keysMap[hash(key)] = i
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (d *Document) StringsForLanguage(lang string) []string {
	return d.Translations[lang]
}

func (d *Document) StringForLanguageAndKey(lang, key string) (string, bool) {
	keyIndex, found := d.keysMap[hash(key)]
	if !found {
		return "", false
	}
	values, found := d.Translations[lang]
	if !found || len(values) < keyIndex {
		return "", false
	}
	return values[keyIndex], true
}

func (d *Document) set(language string, values []string) {
	if d.Translations == nil {
		d.Translations = make(map[string][]string, len(d.LanguageNames))
	}
	d.Translations[language] = values
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
		stringsPath := filepath.Join(basepath, dir, lang+".lproj", file+".strings")
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
		err = ioutil.WriteFile(stringsPath, []byte(content.String()), exportFileMode)
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

type androidStrings struct {
	XMLName xml.Name `xml:"resources"`
	Items   []androidStringsItem
}

type androidStringsItem struct {
	XMLName xml.Name `xml:"string"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",innerxml"`
}

func (doc *Document) exportAsAndroidStrings(basepath string) error {
	for _, lang := range doc.LanguageNames {
		subdir := "values"
		if lang != "" {
			subdir += "-" + lang
		}
		stringsPath := filepath.Join(basepath, doc.Path, subdir, "strings.xml")
		translations := doc.stringMapForLanguage(lang)
		if len(translations) == 0 {
			println("translations for", lang, "is empty.")
			continue
		}
		println("writing document", doc.Name, "(", len(doc.KeyNames), "keys) to\n", stringsPath)
		dir := path.Dir(stringsPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0700)
		}
		var it androidStrings
		for _, key := range doc.KeyNames {
			text := translations[key]
			it.Items = append(it.Items, androidStringsItem{Name: key, Value: text})
		}
		bs, err := xml.MarshalIndent(it, "", "\t")
		if err != nil {
			return err
		}
		bs = append([]byte(xml.Header), bs...)
		err = ioutil.WriteFile(stringsPath, bs, exportFileMode)
		if err != nil {
			return err
		}
	}
	return nil
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
	bs, err := json.MarshalIndent(map[string]interface{}{
		"keys": doc.KeyNames,
		"text": translations,
	}, "", "")
	if err != nil {
		return err
	}
	exportPath := filepath.Join(dir, file+".json")
	println("writing document", doc.Name, "(", len(doc.KeyNames), "keys) to\n", exportPath)
	return ioutil.WriteFile(exportPath, bs, exportFileMode)
}

func (doc *Document) exportAsARB(basepath string) error {
	for _, lang := range doc.LanguageNames {
		dir, file := doc.pathComponents()
		docPath := filepath.Join(basepath, dir, file+"_"+lang+".arb")
		translations := doc.stringMapForLanguage(lang)
		if len(translations) == 0 {
			println("translations for", lang, "is empty.")
			continue
		}
		println("writing document", doc.Name, "(", len(doc.KeyNames), "keys) to\n", docPath)
		docDir := path.Dir(docPath)
		if _, err := os.Stat(docDir); os.IsNotExist(err) {
			os.MkdirAll(docDir, 0700)
		}
		var err error
		content := map[string]interface{}{
			"@@locale": lang,
			// "@@last_modified": time.Now().UTC().Format(time.RFC3339),
		}
		for _, key := range doc.KeyNames {
			text := translations[key]
			if text != "" {
				content[key] = text
			}
		}
		bs, err := json.MarshalIndent(content, "", "")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(docPath, bs, exportFileMode)
		if err != nil {
			return err
		}
		listPath := filepath.Join(basepath, dir, file+".arb")
		bs, err = json.Marshal(map[string]interface{}{
			"locales": doc.LanguageNames,
		})
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(listPath, bs, exportFileMode)
		if err != nil {
			return err
		}
	}
	return nil
}
