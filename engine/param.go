package engine

import (
	"flag"
	"fmt"
	"os"

	"github.com/CaiJinKen/gostruct/source"
)

func init() {
	flag.StringVar(&defaultParam.inputFile, "i", "", "input sql file")
	flag.StringVar(&defaultParam.outputFile, "o", "", "output file")
	flag.BoolVar(&defaultParam.echo, "e", true, "echo result")
	flag.BoolVar(&defaultParam.useGormTag, "g", false, "use gorm tag")
	flag.BoolVar(&defaultParam.useJsonTag, "j", true, "use json tag")
	flag.BoolVar(&defaultParam.sortField, "s", false, "sort fields by ASCII")
	flag.StringVar(&defaultParam.pkgName, "p", "models", "package name")
	flag.StringVar(&defaultParam.dsn, "d", "", "mysql dsn, format: 'user:password@tcp(host:port)/db_name'")
	flag.StringVar(&defaultParam.tableName, "t", "", "table name")
}

type param struct {
	inputFile  string // input sql file
	outputFile string // result file
	echo       bool   // echo result
	useGormTag bool
	useJsonTag bool
	sortField  bool   // sort field
	pkgName    string // output file package name
	dsn        string // sql dsn
	tableName  string // table name
}

var defaultParam param

func getParam() *param {
	return &defaultParam
}

func (p *param) GetSource() source.Source {
	if p.inputFile != "" {
		return &source.FileSource{
			FilePath: p.inputFile,
		}
	}
	if p.dsn != "" && p.tableName != "" {
		return &source.DbSource{
			DSN:       p.dsn,
			TableName: p.tableName,
		}
	}
	fmt.Println("need one table source (input file or dsn) at lest.")
	os.Exit(1)
	return nil
}
