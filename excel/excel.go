package excel

import (
	"github.com/tealeg/xlsx"
	"strings"
)

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
		keyCount := sheet.MaxRow - 1
		keys := make([]string, keyCount)
		for i := 0; i < sheet.MaxRow; i += 1 {
			if i > 0 {
				keys[i - 1] = sheet.Cell(i, 0).Value
			}
		}
		doc.SetKeys(keys)

		for ci := 1; ci < sheet.MaxCol; ci += 1 {
			lang := doc.LanguageNames[ci - 1]
			values := make([]string, keyCount)
			for ri := 1; ri < sheet.MaxRow; ri += 1 {
				cell := sheet.Cell(ri, ci)
				values[ri - 1] = cell.Value
			}
			doc.set(lang, values)
		}
		docs.Documents = append(docs.Documents, doc)
	}
	return docs, nil
}
