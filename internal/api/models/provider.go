package models

import (
	"auric/internal/providers"
	"auric/internal/providers/consul"
)

var Provider providers.AuricProvider

func InitProvider(providerType string, config *providers.ProviderConfig) {
	if providerType == "consul" {
		p := &consul.ConsulClient{}

		//TODO read from json file
		err := p.Init(*config)

		if err != nil {
			panic(err)
		}
		Provider = p
	}

	/*
		if providerType == "test" {
			p := &test.TestClient{}

			//TODO read from json file
			err := p.Init(*config)

			if err != nil {
				panic(err)
			}
			Provider = p
		}*/
}
