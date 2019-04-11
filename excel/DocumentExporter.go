package excel

import (
	"encoding/json"
	"io/ioutil"
	"os"
	path2 "path"
	strings2 "strings"
)

func ExportIOSStrings(docs *Documents, basepath string) error {
	for _, doc := range docs.Documents {
		for _, lang := range doc.LanguageNames {
			path := basepath + "/" + doc.Path + "/" + lang + ".lproj/" + doc.File + ".strings"
			if err := exportIOSForLanguageAtPath(doc, lang, path); err != nil {
				return err
			}
		}
	}
	return nil
}

func exportIOSForLanguageAtPath(doc *Document, lang, path string) error {
	strings := doc.StringsForLanguage(lang)
	if len(strings) == 0 {
		return nil
	}
	println("writing document", doc.Name, "(", len(doc.KeyNames), "keys) to\n", path)
	dir := path2.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(path, 0666)
	}
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	for _, key := range doc.KeyNames {
		text := strings[key]
		if strings2.Contains(key, " ") {
			key = "\"" + key + "\""
		}
		text = strings2.Replace(text, "\"", "\\\"", -1)
		text = strings2.Replace(text, "\n", "\\n", -1)
		line := key + "=\"" + text + "\";\n"
		_, err = fp.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExportI18n(docs *Documents, basepath string) error {
	for _, doc := range docs.Documents {
		for _, lang := range doc.LanguageNames {
			path := basepath + "/" + doc.Path + "/" + doc.File + "_" + lang + ".json"
			if err := exportI18nAtPathForLanguage(doc, path, lang, "en"); err != nil {
				return err
			}
		}
	}
	return nil
}

func exportI18nAtPathForLanguage(doc *Document, path, lang, defLang string) error {
	strings := doc.StringsForLanguage(lang)
	if len(strings) == 0 {
		return nil
	}
	keyCount := len(doc.KeyNames)
	println("writing document", doc.Name, "(", keyCount, "keys) to\n", path)
	dir := path2.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(path, 0666)
	}
	m := make(map[string]string, keyCount)
	defStrings := doc.StringsForLanguage(defLang)
	for _, key := range doc.KeyNames {
		text := strings[key]
		if len(text) == 0 {
			text = defStrings[key]
		}
		m[key] = text
	}
	bs, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bs, 0666)
}
