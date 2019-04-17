package excel

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func ExportDocumentsAsIOSStrings(docs *Documents, basepath string) error {
	for _, doc := range docs.Documents {
		for _, lang := range doc.LanguageNames {
			stringsPath := basepath + "/" + doc.Path + "/" + lang + ".lproj/" + doc.File + ".strings"
			translations := StringMapInDocumentForLanguage(doc, lang)
			if len(translations) == 0 {
				println("translations for", lang, "is empty.")
				continue
			}
			println("writing document", doc.Name, "(", len(doc.KeyNames), "keys) to\n", stringsPath)
			dir := path.Dir(stringsPath)
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
	}
	return nil
}

func StringMapInDocumentForLanguage(d *Document, lang string) map[string]string {
	m := make(map[string]string, len(d.KeyNames))
	for _, key := range d.KeyNames {
		v, found := d.StringForLanguageAndKey(lang, key)
		if found {
			m[key] = v
		}
	}
	return m
}
