package main

import "os"
import "github.com/shawnclovie/Localizable.excel.go/excel"

func main() {
	args := os.Args[1:]
	argc := len(args)
	if argc < 2 {
		print("usage: ", os.Args[0], " <export_ios|export_i18n> <xlsx_file_name on current directory>")
		return
	}
	projDir, err := os.Getwd()
	panicIfNotNull(err)

	action := args[0]
	filename := args[1]
	excelPath := projDir + "/" + filename

	docs, err := excel.LoadExcelFromFile(excelPath)
	panicIfNotNull(err)
	switch action {
	case "export_ios":
		panicIfNotNull(excel.ExportIOSStrings(docs, projDir))
	case "export_i18n":
		panicIfNotNull(excel.ExportI18n(docs, projDir))
	default:
		panic("undefined operation")
	}
}

func panicIfNotNull(err error) {
	if err != nil {
		panic(err)
	}
}
