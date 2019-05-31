package src

import (
	"fmt"
	"os"
)

func Run() {
	if inputFile == nil || *inputFile == "" {
		fmt.Println("no input file.")
		os.Exit(-1)
	}

	file, err := os.OpenFile(*inputFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("open file err: ", err.Error())
		os.Exit(-1)
	}
	defer file.Close()

	table := parseFile(file)
	buf := parseTable(table)
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
