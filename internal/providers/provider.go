package providers

type AuricProvider interface {
	Init(config ProviderConfig) error
	Store(path string, value []byte) error
	Retrieve(path string) ([]byte, error)
	Delete(path string) error
}

type ProviderConfig struct {
	ConsulConfig struct {
		Address    string `json:"address"`
		Port       string `json:"port"`
		Datacenter string `json:"datacenter"`
	} `json:"consul_config"`
}
