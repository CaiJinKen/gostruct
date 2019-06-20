package src

import (
	"fmt"
	"io/ioutil"
	"os"
)

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
			fmt.Printf("read fiile err: %+v\n", err)
		}
		data = bts
	}
	return
}

func getFromDB() (data []byte) {
	if dsn != nil && *dsn != "" {
		db := getMysqlDB(*dsn)
		if db == nil {
			fmt.Printf("get db connection err.\n")
			return
		}
		defer db.Close()
		data = getTable(db)
	}
	return
}
