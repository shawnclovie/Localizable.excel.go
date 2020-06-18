A [go](http://www.golang.org) (or 'golang' for search engine friendliness) implementation of localize strings render tool.

Edit text in Excel and export to iOS strings or json.

## Structure of Excel
Each work sheet is a Localizabe.strings group:

* Cell[0,0] should be -
```
path={PathToSaving}/{LocalizableFileName}
format=json|ios_strings
```
Each argument separate with & or new line.

* Other cells on row[0] should be language code, e.g. en, zh-Hans, it.

* Other cells on column[0] should be strings key, what you typed in code - NSLocalizableString("StringKey", comment: nil).

* Cells on column[1] is primary language, it should be fully filled.

|path=Resources/Localizable&format=ios_strings|en|ja|
|-|-|-|
|home.title|Text Editor|テキストエディタ|
|home.footer|Copyright|

## Export
* Command
```
$ go run main.go export excel.xlsx
```

* format=ios_strings
```
Resources/en.lproj/Localizable.strings
home.title="Text Editor"
home.footer="Copyright"

Resources/ja.lproj/Localizable.strings
home.title="テキストエディタ"
```

* format=json
```
Resources/Localizable_en.json
{"home.title":"Text Editor","home.footer":"Copyright"}
Resources/Localizable_ja.json
{"home.title":"テキストエディタ","home.footer":"Copyright"}
```

## Dependence

* github.com/tealeg/xlsx
