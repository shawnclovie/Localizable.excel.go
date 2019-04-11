package excel

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func ExportDocumentsAsIOSStrings(docs *Documents, basepath string) error {
	for _, doc := range docs.Documents {
		for _, lang := range doc.LanguageNames {
			stringsPath := basepath + "/" + doc.Path + "/" + lang + ".lproj/" + doc.File + ".strings"
			translations := doc.StringMapForLanguage(lang)
			if len(translations) == 0 {
				continue
			}
			println("writing document", doc.Name, "(", len(doc.KeyNames), "keys) to\n", stringsPath)
			dir := path.Dir(stringsPath)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.MkdirAll(dir, 0700)
			}
			fp, err := os.Open(stringsPath)
			if err != nil {
				return err
			}
			for _, key := range doc.KeyNames {
				text := translations[key]
				if strings.Contains(key, " ") {
					key = "\"" + key + "\""
				}
				text = strings.Replace(text, "\"", "\\\"", -1)
				text = strings.Replace(text, "\n", "\\n", -1)
				line := key + "=\"" + text + "\";\n"
				_, err = fp.WriteString(line)
				if err != nil {
					break
				}
			}
			fp.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
