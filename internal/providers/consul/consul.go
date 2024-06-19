package consul

import (
	"fmt"
	"os"

	capi "github.com/hashicorp/consul/api"
)

var consul_address = "127.0.0.1:8500"
var consul_datacenter = "dc1"

func (c Consul) Init() error {
	client, err := newConsulClient()
	c.client = client
	fmt.Println("Initializing consul client")
	return err
}

func (c Consul) Store(path string, value []byte) error {
	err := CreateConsulKVPair(path, value)
	return err
}

func (c Consul) Retrieve(path string) ([]byte, error) {
	pair, err := GetConsulKVPair(path)
	return pair.Value, err
}

func (c Consul) Delete(path string) error {
	return nil
}

func newConsulClient() (*capi.Client, error) {
	token := os.Getenv("CONSUL_HTTP_TOKEN")

	client, err := capi.NewClient(&capi.Config{
		Address:    consul_address,
		Scheme:     "http",
		Datacenter: consul_datacenter,
		Token:      token,
	})
	if err != nil {
		panic(err)
	}
	return client, err
}

func CreateConsulKVPair(key string, value []byte) error {
	consul_client, err := newConsulClient()
	if err != nil {
		panic(err)
	}

	kv := consul_client.KV()

	pair := &capi.KVPair{
		Key:   key,
		Value: value,
	}
	_, err = kv.Put(pair, nil)
	return err
}

func GetConsulKVPair(key string) (*capi.KVPair, error) {
	consul_client, err := newConsulClient()
	if err != nil {
		panic(err)
	}

	kv := consul_client.KV()
	pair, _, err := kv.Get(key, nil)

	return pair, err
}