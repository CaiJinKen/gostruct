package src

import "flag"

var (
	linLin        = []byte{lin, lin}
	slashStar     = []byte{slash, star}
	starSlash     = []byte{star, slash}
	set           = []byte("SET")
	drop          = []byte("DROP")
	space         = []byte{' '}
	point         = []byte{'.'}
	unsigned      = []byte("UNSIGNED")
	comment       = []byte("COMMENT")
	notNull       = []byte("NOT NULL")
	autoIncrement = []byte("AUTO_INCREMENT")
	dft           = []byte("DEFAULT")
)

var (
	tmpFile          = "gostruct_temp_struct.go"
	tableNameFuncStr = `
func (*%s) TableName() string {
	return "%s"
}
`
)

var (
	star  byte = '*'
	slash byte = '/'
	lin   byte = '-'
)

var (
	inputFile  = flag.String("i", "", "input sql file")
	outputFile = flag.String("o", "", "output file")
	echo       = flag.Bool("e", true, "echo result")
	gormTag    = flag.Bool("g", false, "use gorm tag")
	jsonTag    = flag.Bool("j", true, "use json tag")
	sortField  = flag.Bool("s", false, "sort fields by ASCII")
	pkg        = flag.String("p", "models", "package name")
	dsn        = flag.String("d", "", "mysql dsn, format: 'user:password@tcp(host:port)/db_name'")
	tableName  = flag.String("t", "", "table name")
)
