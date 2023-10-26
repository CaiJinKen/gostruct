package table

import (
	"bytes"
	"os"
	"os/exec"
	"sync"

	"github.com/CaiJinKen/gostruct/handler"
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
	m.once.Do(func() {
		cmd := exec.Command("gofmt")
		cmd.Stdin = m.reader
		cmd.Stdout = m.writer
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			handler.PrintErrAndExit(err)
			return
		}
		m.result = m.writer.Bytes()
		m.writer = nil
		m.reader = nil
	})

}
