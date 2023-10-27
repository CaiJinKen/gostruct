package handler

import (
	"fmt"
	"os"
)

func PrintErrAndExit(err error) {
	fmt.Printf("gostruct parse err %+v\n", err)
	os.Exit(1)
}
