package excel

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ExportDocumentsAsJSON(docs *Documents, basepath string) error {
	for _, doc := range docs.Documents {
		dir := basepath + "/" + doc.Path
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
		exportPath := dir + "/" + doc.File + ".json"
		println("writing document", doc.Name, "(", len(doc.KeyNames), "keys) to\n", exportPath)
		err = ioutil.WriteFile(exportPath, bs, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
