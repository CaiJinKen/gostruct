package src

import (
	"fmt"
	"os"
	"path/filepath"
)

func Run() {
	if (inputFile == nil || *inputFile == "") && (dsn == nil || *dsn == "") {
		fmt.Println("need one table source (input file or dsn) at lest.")
		os.Exit(1)
	}

	table := parseTable(getTableBytes())
	buf := marshalTable(&table)
	buf = writeTmpFile(buf)
	if outputFile != nil && *outputFile != "" {
		absFile()
		os.MkdirAll(filepath.Dir(*outputFile), os.ModePerm)
		os.Rename(tmpFile, *outputFile)
	}
	defer os.Remove(tmpFile)

	if *echo {
		fmt.Println(string(buf))
	}
}
