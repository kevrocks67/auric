package consul

import (
	"auric/internal/providers"
	"errors"
	"fmt"
	"os"

	capi "github.com/hashicorp/consul/api"
)

func (c *ConsulClient) Init(config providers.ProviderConfig) error {
	client, err := NewConsulClient(config)
	c.Client = client
	fmt.Println("Initializing consul client")
	return err
}

func (c *ConsulClient) Store(path string, value []byte) error {
	err := CreateConsulKVPair(c, path, value)
	return err
}

func (c *ConsulClient) Retrieve(path string) ([]byte, error) {
	pair, err := GetConsulKVPair(c, path)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, errors.New(fmt.Sprintf("No key pair found for path: %s", path))
	}
	return pair.Value, nil
}

func (c *ConsulClient) Delete(path string, isPrefix bool) error {
	var err error

	if isPrefix {
		err = DeleteConsulKVPath(c, path)
	} else {
		err = DeleteConsulKVPair(c, path)
	}
	return err
}

func (c *ConsulClient) List(prefix string) ([][]byte, error) {
	pairs, err := ListConsulKVPath(c, prefix)
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

func NewConsulClient(config providers.ProviderConfig) (*capi.Client, error) {
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
	kv := c.Client.KV()

	pair := &capi.KVPair{
		Key:   key,
		Value: value,
	}

	_, err := kv.Put(pair, nil)
	return err
}

func GetConsulKVPair(c *ConsulClient, key string) (*capi.KVPair, error) {
	kv := c.Client.KV()
	pair, _, err := kv.Get(key, nil)

	return pair, err
}

func ListConsulKVPath(c *ConsulClient, prefix string) ([][]byte, error) {
	pairs, _, err := c.Client.KV().List(prefix, nil)
	var pairList [][]byte
	for _, pair := range pairs {
		pairList = append(pairList, pair.Value)
	}
	return pairList, err
}

func DeleteConsulKVPair(c *ConsulClient, path string) error {
	kv := c.Client.KV()

	_, err := kv.Delete(path, nil)
	return err
}

func DeleteConsulKVPath(c *ConsulClient, prefix string) error {
	kv := c.Client.KV()

	_, err := kv.DeleteTree(prefix, nil)
	return err
}
