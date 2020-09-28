package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/shawnclovie/Localizable.excel.go/excel"
)

const (
	fromNew = "new"
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
 new <json|xlsx> <file> <languages, separate with ","> <doc names>
Export:
 <json|xlsx> export <input file> <export dir>
Convert:
 xlsx json <input file> [output file]
 json xlsx <input file> [output file]
`)
		return
	}
	fmt.Println("args:", args)

	from := strings.TrimSpace(args[0])
	action := strings.TrimSpace(args[1])
	if from == action {
		println("1st and 2nd arguments should not equal.")
		os.Exit(1)
	}
	filename := args[2]
	inputPath := filename
	if !strings.HasPrefix(filename, "/") {
		cwd, err := os.Getwd()
		panicIfNotNull(err)
		inputPath = cwd + "/" + filename
	}

	var docs excel.Documents
	var err error
	var savePath string
	switch from {
	case formatXLSX:
		docs, err = excel.LoadDocumentsFromExcelFile(inputPath)
	case formatJSON, formatYAML:
		docs, err = excel.LoadDocumentsFromFile(inputPath)
	case fromNew:
		langs := strings.Split(args[3], ",")
		docNames := strings.Split(args[4], ",")
		docs = excel.NewDocument(langs, docNames)
		savePath = inputPath
	default:
		err = errors.New("undefined operation")
	}
	panicIfNotNull(err)
	if action == actionExport {
		exportDir := args[3]
		println(from, action, len(docs.Documents), "Sheet(s)", "to", exportDir)
		panicIfNotNull(docs.Export(exportDir))
		return
	}
	if savePath == "" {
		if len(args) > 3 {
			savePath = args[3]
		} else {
			ext := filepath.Ext(inputPath)
			savePath = inputPath[:len(inputPath)-len(ext)] + "." + action
		}
	}
	print("convert (", from, ")\n", inputPath, "\nto (", action, ")\n", savePath, "\n")
	var bs []byte
	switch action {
	case formatJSON:
		bs, err = docs.ToJSONData()
	case formatYAML:
		bs, err = docs.ToYAMLData()
	case formatXLSX:
		bs, err = docs.ToExcelData()
	default:
		fmt.Println(docs)
		err = errors.New("undefined operation")
	}
	panicIfNotNull(err)

	err = ioutil.WriteFile(savePath, bs, 0644)
	panicIfNotNull(err)
	println("done")
}

func panicIfNotNull(err error) {
	if err != nil {
		println(err.Error())
		panic(err)
	}
}
