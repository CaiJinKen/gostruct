package src

import "flag"

var (
	tmpFile            = "temp_struct.go"
	star          byte = '*'
	slash         byte = '/'
	lin           byte = '-'
	linLin             = []byte{lin, lin}
	slashStar          = []byte{slash, star}
	starSlash          = []byte{star, slash}
	set                = []byte("SET")
	drop               = []byte("DROP")
	space              = []byte{' '}
	point              = []byte{'.'}
	unsigned           = []byte("UNSIGNED")
	comment            = []byte("COMMENT")
	notNull            = []byte("NOT NULL")
	autoIncrement      = []byte("AUTO_INCREMENT")
	dft                = []byte("DEFAULT")

	tableNameFuncStr = `
func (*%s) TableName() string {
	return "%s"
}
`
)

var inputFile = flag.String("i", "", "input sql file")
var outputFile = flag.String("o", "", "output file")
var echo = flag.Bool("e", true, "echo result")
var gormTag = flag.Bool("g", false, "use gorm tag")
var jsonTag = flag.Bool("j", true, "use json tag")
var sortField = flag.Bool("s", false, "sort fields by ASCII")
var pkg = flag.String("p", "models", "package name")
var dsn = flag.String("d", "", "mysql dsn, format: user:password@tcp(host:port)/db_name")
var tableName = flag.String("t", "", "table name")
