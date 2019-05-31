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
)

var inputFile = flag.String("i", "", "input sql file")
var outputFile = flag.String("o", "", "output file")
var echo = flag.Bool("e", true, "echo result")
var gorm = flag.Bool("g", false, "use gorm tag")
var sortField = flag.Bool("s", false, "sort fields by ASCII")
var pkg = flag.String("p", "main", "package name")
