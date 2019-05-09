package src

import "flag"

var (
	tmpFile        = "temp_struct.go"
	star      byte = '*'
	slash     byte = '/'
	lin       byte = '-'
	linlin         = []byte{lin, lin}
	slashStar      = []byte{slash, star}
	starSlash      = []byte{star, slash}
	set            = []byte("SET")
	drop           = []byte("DROP")
	space          = []byte{' '}
	point          = []byte{'.'}
	unsigns        = []byte("UNSIGNED")
	comment        = []byte("COMMENT")
)

var inputFile = flag.String("i", "", "input sql file")
var outputFile = flag.String("o", "", "output file")
var print = flag.Bool("p", true, "if print result")
