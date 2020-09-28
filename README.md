A [go](http://www.golang.org) (or 'golang' for search engine friendliness) implementation of localize strings render tool.

Edit text in Excel or formatted JSON/YAML and export to:
- iOS **strings**
- Android **string.xml**
- Flutter **arb** json for intl_localization

## Structure of Formatted JSON or YAML
```yaml
- format: ios_string
  language_names: [en,ja,zh,zh-Hant]
  name: iOS
  path: Resources/Localizable
  translations:
    - [home.title,Text Editor,テキストエディタ,文字编辑,文字編輯]
	- [home.footer,Copyright]
- format: android
  path: app/src/main/res
  name: Android
  language_names: ["",en]
  translations:
    - [app_name,了不起的应用程序,Amazing App]
- format: arb
  ...
```

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

|path=Resources/Localizable&format=ios|en|ja|zh|zh-Hant-HK
|-|-|-|
|home.title|Text Editor|テキストエディタ|文字编辑|文字編輯|
|home.footer|Copyright|

|path=app/src/main/res&format=android||en|
|-|-|-|
|app_name|了不起的应用程序|Amazing App|

## Help
```
$ go run main.go
```

## Convert
```
$ go run main.go <json|yaml|xlsx> <json|yaml|xlsx> <input file>
```

## Export
* Command
```
$ go run main.go <json|yaml|xlsx> export <input file> <export dir>
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
	* Resources/Localizabe.arb
	```json
	{"locales":["en","ja"]}
	```
	* Resources/Localizabe_en.arb
	```json
	{"@@locale":"en","home.title":"Text Editor","home.footer":"Copyright"}
	```
	* Resources/Localizabe_ja.arb
	```json
	{"@@locale":"ja","home.title":"テキストエディタ"}
	```

## Dependence

* github.com/tealeg/xlsx
* gopkg.in/yaml.v3
