package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/shawnclovie/Localizable.excel.go/excel"
	"github.com/shawnclovie/Localizable.excel.go/utility"
)

const (
	fromNew    = "new"
	formatXLSX = "xlsx"
	formatJSON = "json"
	formatYAML = "yaml"

	actionExport = "export"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		println("usage: ")
		println(`
New file:
 new <file> <languages, separate with ","> <doc names>
- file should has extension json | yaml | xlsx.

Export:
 export <input file> <export dir>

Convert input file to target format:
 convert <input file> <output file>
`)
		return
	}
	fmt.Println("args:", args)

	var act action
	switch strings.ToLower(strings.TrimSpace(args[0])) {
	case "new":
		langs := strings.Split(args[2], ",")
		docNames := strings.Split(args[3], ",")
		act = newDocsAction{
			targetFile: args[1],
			docs:       excel.NewDocument(langs, docNames),
		}
	case "export":
		src := args[1]
		act = exportAction{
			sourceFile: src,
			exportDir:  args[2],
			docs:       readDocuments(src),
		}
	case "convert":
		src := args[1]
		act = convertAction{
			sourceFile: src,
			targetFile: args[2],
			docs:       readDocuments(src),
		}
	}
	act.invoke()
	println("done")
}

type action interface {
	invoke()
}

type newDocsAction struct {
	targetFile string

	docs excel.Documents
}

func (act newDocsAction) invoke() {
	convertAction{
		targetFile: act.targetFile,
		docs:       act.docs,
	}.invoke()
}

type exportAction struct {
	sourceFile string
	exportDir  string

	docs excel.Documents
}

func (act exportAction) invoke() {
	println("export", len(act.docs.Documents), "Sheet(s)", "to", act.exportDir)
	utility.PanicIfNotNull(act.docs.Export(act.exportDir))
}

type convertAction struct {
	sourceFile string
	targetFile string

	docs excel.Documents
}

func (act convertAction) invoke() {
	writeDocumentsToFile(act.docs, act.targetFile)
}

func readDocuments(src string) (docs excel.Documents) {
	if _, err := os.Stat(src); err != nil {
		utility.PanicIfNotNull(err)
	}
	inputExt := filepath.Ext(src)
	if inputExt == "" {
		panic(fmt.Errorf("input file should has suffix: %s", src))
	}
	inputExt = strings.ToLower(inputExt[1:])
	var err error
	switch inputExt {
	case formatXLSX:
		docs, err = excel.LoadDocumentsFromExcelFile(src)
	case formatJSON, formatYAML:
		docs, err = excel.LoadDocumentsFromFile(src)
	default:
		err = errors.New("undefined operation")
	}
	utility.PanicIfNotNull(err)
	return
}

func writeDocumentsToFile(docs excel.Documents, tar string) {
	ext := filepath.Ext(tar)
	if ext == "" {
		panic(fmt.Errorf("target file (%v) should has valid extension", tar))
	}
	ext = ext[1:]
	var bs []byte
	var err error
	switch ext {
	case formatJSON:
		bs, err = docs.ToJSONData()
	case formatYAML:
		bs, err = docs.ToYAMLData()
	case formatXLSX:
		bs, err = docs.ToExcelData()
	default:
		err = errors.New("undefined operation")
	}
	utility.PanicIfNotNull(err)
	err = ioutil.WriteFile(tar, bs, 0644)
	utility.PanicIfNotNull(err)
}
