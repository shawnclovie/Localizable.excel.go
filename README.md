A [go](http://www.golang.org) (or 'golang' for search engine friendliness) implementation of localize strings render tool.

Edit text in Excel and export to:
- iOS **strings**
- Android **string.xml**
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

|path=Resources/Localizable&format=ios|en|ja|
|-|-|-|
|home.title|Text Editor|テキストエディタ|
|home.footer|Copyright|

|path=app/src/main/res&format=android||en|
|-|-|-|
|app_name|了不起的应用程序|Amazing App|

## Export
* Command
```
$ go run main.go export excel.xlsx
```

* format=ios
	* Resources/en.lproj/Localizable.strings
	```
	home.title="Text Editor"
	home.footer="Copyright"
	```
	* Resources/ja.lproj/Localizable.strings
	```
	home.title="テキストエディタ"
	```
* format=android
	* app/src/main/res/value/strings.xml
	```xml
	<?xml version="1.0" encoding="UTF-8"?>
	<resources>
		<string name="app_name">了不起的应用程序</string>
	</resources>
	```
	* app/src/main/res/value-en/strings.xml
	```xml
	<?xml version="1.0" encoding="UTF-8"?>
	<resources>
		<string name="app_name">Amazing App</string>
	</resources>
	```

* format=json
	* Resources/Localizable.json
	```json
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
	```json
	{"home.title":"Text Editor","home.footer":"Copyright"}
	```
	* Resources/Localizabe_ja.arb
	```json
	{"home.title":"テキストエディタ"}
	```

## Dependence

* github.com/tealeg/xlsx
