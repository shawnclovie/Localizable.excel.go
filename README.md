A [go](http://www.golang.org) (or 'golang' for search engine friendliness) implementation of localize strings render tool.

Edit text in Excel and export to iOS strings or json.

## Structure of Excel
Each work sheet is a Localizabe.strings group:

* Cell[0,0] should be "PathToSaving\nLocalizableFileName".
* Other cells on row[0] should be language code, e.g. en, zh-Hans, it.
* Other cells on column[0] should be strings key, what you typed in code - NSLocalizableString("StringKey", comment: nil).
* Cells on column[1] is primary language, it should be fully filled.

		Resources
		Localizable   en           ja
		home.title    Text Editor  テキストエディタ
		home.footer   Copyright

It would rendered to -

* $go run main.go export_ios excel.xlsx
```
Resources/en.lproj/Localizable.strings
Resources/ja.lproj/Localizable.strings
```
* $go run main.go export_i18n excel.xlsx
```
Resources/Localizable_en.json
Resources/Localizable_ja.json
```
