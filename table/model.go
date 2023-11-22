package table

import (
	"bytes"
	"go/format"
	"sync"
)

type model struct {
	reader, writer *bytes.Buffer
	result         []byte
	once           sync.Once
}

func newModel() *model {
	return &model{
		reader: &bytes.Buffer{},
		writer: &bytes.Buffer{},
	}
}

func (m *model) format() {
	if len(m.result) > 0 {
		return
	}
	if m.reader == nil || m.reader.Len() == 0 {
		return
	}

	m.result, _ = format.Source(m.reader.Bytes())

	// m.once.Do(func() {
	// 	cmd := exec.Command("gofmt")
	// 	cmd.Stdin = m.reader
	// 	cmd.Stdout = m.writer
	// 	cmd.Stderr = os.Stderr
	//
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		handler.PrintErrAndExit(err)
	// 		return
	// 	}
	// 	m.result = m.writer.Bytes()
	// 	m.writer = nil
	// 	m.reader = nil
	// })

}
