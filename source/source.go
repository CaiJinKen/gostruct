package source

var (
	_ Source = new(FileSource)
	_ Source = new(DbSource)
)

type Source interface {
	GetData() (data []byte, err error)
}
