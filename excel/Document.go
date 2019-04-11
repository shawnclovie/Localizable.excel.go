package excel

import (
	"github.com/tealeg/xlsx"
	"strings"
)

type Documents struct {
	Path string
	Documents []*Document
}

type Document struct {
	Name string
	Path string
	File string
	LanguageNames []string
	KeyNames []string
	strings map[string]map[string]string
}

func (d *Document) StringsForLanguage(lang string) map[string]string {
	return d.strings[lang]
}

func (d *Document) StringForLanguageAndKey(lang, key string) (string, bool) {
	str, found := d.strings[lang][key]
	return str, found
}

func (d *Document) set(language, key, translation string) {
	if d.strings == nil {
		d.strings = map[string]map[string]string{}
	}
	if d.strings[language] == nil {
		d.strings[language] = map[string]string{}
	}
	d.strings[language][key] = translation
}

func LoadExcelFromFile(path string) (*Documents, error) {
	xls, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}
	docs := &Documents{Path: path}
	for _, sheet := range xls.Sheets {
		if sheet.MaxCol <= 1 || sheet.MaxRow <= 1 {
			continue
		}
		names := strings.Split(sheet.Cell(0, 0).Value, "\n")
		if len(names) < 2 {
			continue
		}
		row0 := sheet.Rows[0]
		doc := &Document{Name: sheet.Name, Path: names[0], File: names[1]}
		doc.LanguageNames = make([]string, len(row0.Cells) - 1)
		for index, cell := range row0.Cells {
			if index > 0 {
				doc.LanguageNames[index - 1] = cell.Value
			}
		}
		doc.KeyNames = make([]string, sheet.MaxRow - 1)
		for i := 0; i < sheet.MaxRow; i += 1 {
			if i > 0 {
				doc.KeyNames[i - 1] = sheet.Cell(i, 0).Value
			}
		}
		for ri := 1; ri < sheet.MaxRow; ri += 1 {
			for ci := 1; ci < sheet.MaxCol; ci += 1 {
				cell := sheet.Cell(ri, ci)
				doc.set(doc.LanguageNames[ci - 1], doc.KeyNames[ri - 1], cell.Value)
			}
		}
		docs.Documents = append(docs.Documents, doc)
	}
	return docs, nil
}
