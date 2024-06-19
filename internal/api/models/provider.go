package models

import (
	"auric/internal/providers"
	"auric/internal/providers/consul"
)

var Provider providers.AuricProvider

func InitProvider(providerType string) {
	if providerType == "consul" {
		p := &consul.ConsulClient{}

		//TODO read from json file
		config := providers.ProviderConfig{
			ConsulConfig: struct {
				Address    string `json:"address"`
				Port       string `json:"port"`
				Datacenter string `json:"datacenter"`
			}{
				Address:    "127.0.0.1",
				Port:       "8500",
				Datacenter: "dc1",
			},
		}
		err := p.Init(config)

		if err != nil {
			panic(err)
		}
		Provider = p
	}
}
