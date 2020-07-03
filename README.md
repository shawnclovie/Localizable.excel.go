A [go](http://www.golang.org) (or 'golang' for search engine friendliness) implementation of localize strings render tool.

Edit text in Excel and export to:
- iOS **strings**
- Key sorted **json**
- Flutter **arb** json for intl_localization

## Structure of Excel
Each work sheet is a Localizabe.strings group:

* Cell[0,0] should be -
```
path={PathToSaving}/{LocalizableFileName}
format=ios_strings|json|arb
```
Each argument separate with & or new line.

* Other cells on row[0] should be language code, e.g. en, zh-Hans, it.

* Other cells on column[0] should be strings key, what you typed in code - NSLocalizableString("StringKey", comment: nil).

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
	* Resources/en.lproj/Localizable.strings
	```
	home.title="Text Editor"
	home.footer="Copyright"
	```
	* Resources/ja.lproj/Localizable.strings
	```
	home.title="テキストエディタ"
	```

* format=json
	* Resources/Localizable.json
	```
	{
		"keys":["home.title","home.footer"],
		"text":{
			"en":["Text Editor","Copyright"],
			"ja":["テキストエディタ",""]
		}
	}
	```
* format=arb
	* Resources/Localizabe_en.arb
	```
	{"home.title":"Text Editor","home.footer":"Copyright"}
	```
	* Resources/Localizabe_ja.arb
	```
	{"home.title":"テキストエディタ"}
	```

## Dependence

* github.com/tealeg/xlsx
