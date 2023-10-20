package main

import (
	"flag"

	"github.com/CaiJinKen/gostruct/engine"
)

func main() {
	flag.Parse()
	engine := engine.New()
	engine.Run()
}
