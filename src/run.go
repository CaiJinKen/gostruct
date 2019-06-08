package src

import (
	"fmt"
	"io/ioutil"
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

func getTableBytes() []byte {
	bts := getFromFile()
	if len(bts) > 0 {
		return bts
	}

	bts = getFromDB()
	if len(bts) > 0 {
		return bts
	}
	os.Exit(3)
	return nil
}

func getFromFile() (data []byte) {
	if inputFile != nil && *inputFile != "" {
		bts, err := ioutil.ReadFile(*inputFile)
		if err != nil {
			fmt.Sprintf("read fiile err: %+v", err)
		}
		data = bts
	}
	return
}

func getFromDB() (data []byte) {
	if dsn != nil && *dsn != "" {
		db := getMysqlDB(*dsn)
		if db == nil {
			fmt.Sprintf("get db connection err.\n")
			return
		}
		defer db.Close()
		data = getTable(db)
	}
	return
}
