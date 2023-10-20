package source

import (
	"os"
)

type FileSource struct {
	FilePath string
}

func (f *FileSource) GetData() (data []byte, err error) {
	if f.FilePath == "" {
		return
	}

	return os.ReadFile(f.FilePath)
}
