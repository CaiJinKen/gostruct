package src

import (
	"fmt"
	"os"
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
		os.Rename(tmpFile, *outputFile)
	} else {
		os.Remove(tmpFile)
	}
	if *echo {
		fmt.Println(string(buf))
	}
}
