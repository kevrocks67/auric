package providers

type AuricProvider interface {
	Init() error
	Store(path string, value []byte) error
	Retrieve(path string) ([]byte, error)
	Delete(path string) error
}
