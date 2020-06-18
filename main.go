package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/shawnclovie/Localizable.excel.go/excel"
)

func main() {
	args := os.Args[1:]
	argc := len(args)
	if argc < 2 {
		print("usage: ", os.Args[0], " <export_ios|export_json> <xlsx_file_name on current directory>")
		return
	}

	action := args[0]
	filename := args[1]
	excelPath := filename
	if !strings.HasPrefix(filename, "/") {
		cwd, err := os.Getwd()
		panicIfNotNull(err)
		excelPath = cwd + "/" + filename
	}

	docs, err := excel.LoadExcelDocumentsFromFile(excelPath)
	panicIfNotNull(err)
	switch action {
	case "export":
		exportDir := filepath.Dir(excelPath)
		println(action, len(docs.Documents), "Sheet(s)", "to", exportDir)
		panicIfNotNull(docs.Export(exportDir))
	default:
		panic("undefined operation")
	}
}

func panicIfNotNull(err error) {
	if err != nil {
		println(err.Error())
		panic(err)
	}
}
