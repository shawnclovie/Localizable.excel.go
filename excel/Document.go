package excel

import "hash/fnv"

type Documents struct {
	Path string
	Documents []*Document
}

type Document struct {
	Name          string
	Path          string
	File          string
	LanguageNames []string
	KeyNames      []string
	keysMap       map[uint32]int
	translations  map[string][]string
}

func (d *Document) SetKeys(keys []string) {
	d.KeyNames = keys
	d.keysMap = make(map[uint32]int, len(keys))
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

func (d *Document) StringMapForLanguage(lang string) map[string]string {
	m := make(map[string]string, len(d.KeyNames))
	for _, key := range d.KeyNames {
		v, found := d.StringForLanguageAndKey(lang, key)
		if found {
			m[key] = v
		}
	}
	return m
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

