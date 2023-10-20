package engine

import (
	"fmt"
	"os"

	"github.com/CaiJinKen/gostruct/table"
)

type Engine struct {
	param *param
}

func New() *Engine {
	return &Engine{
		param: getParam(),
	}
}

func (e *Engine) Run() {
	data, err := e.rawData()
	if err != nil {
		return
	}

	table := e.Parse(data)
	table.Marshal()
	table.Format()

	data = table.Data()
	e.output(data)

}

func (e *Engine) rawData() ([]byte, error) {
	source := e.param.GetSource()
	return source.GetData()
}

func (e *Engine) Parse(slice []byte) (t *table.Table) {
	c := &table.Config{
		UseGormTag: e.param.useGormTag,
		UseJsonTag: e.param.useJsonTag,
		SortField:  e.param.sortField,
		PkgName:    e.param.pkgName,
	}
	t = c.Build()
	t.Parse(slice)
	return t
}

func (e *Engine) output(data []byte) {
	if e.param.echo {
		fmt.Printf("\n%s\n", string(data))
	}
	e.genFile(data)
}

func (e *Engine) genFile(data []byte) {
	if e.param.outputFile == "" {
		return
	}
	f, err := os.Create(e.param.outputFile)
	if err != nil {
		fmt.Println(fmt.Sprintf("create file %s err %v", e.param.outputFile, err))
		os.Exit(-1)
	}
	defer f.Close()
	f.Write(data)
}
