package consul

import (
	"auric/internal/providers"
	"fmt"
	"os"

	capi "github.com/hashicorp/consul/api"
)

func (c *ConsulClient) Init(config providers.ProviderConfig) error {
	client, err := newConsulClient(config)
	c.client = client
	fmt.Println("Initializing consul client")
	return err
}

func (c *ConsulClient) Store(path string, value []byte) error {
	err := CreateConsulKVPair(c, path, value)
	return err
}

func (c *ConsulClient) Retrieve(path string) ([]byte, error) {
	pair, err := GetConsulKVPair(c, path)
	return pair.Value, err
}

func (c *ConsulClient) Delete(path string) error {
	return nil
}

func newConsulClient(config providers.ProviderConfig) (*capi.Client, error) {
	token := os.Getenv("CONSUL_HTTP_TOKEN")

	client, err := capi.NewClient(&capi.Config{
		Address:    fmt.Sprintf("%s:%s", config.ConsulConfig.Address, config.ConsulConfig.Port),
		Scheme:     "http",
		Datacenter: config.ConsulConfig.Datacenter,
		Token:      token,
	})
	if err != nil {
		panic(err)
	}
	return client, err
}

func CreateConsulKVPair(c *ConsulClient, key string, value []byte) error {
	kv := c.client.KV()

	pair := &capi.KVPair{
		Key:   key,
		Value: value,
	}

	_, err := kv.Put(pair, nil)
	return err
}

func GetConsulKVPair(c *ConsulClient, key string) (*capi.KVPair, error) {
	kv := c.client.KV()
	pair, _, err := kv.Get(key, nil)

	return pair, err
}
