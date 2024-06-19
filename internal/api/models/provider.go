package models

import (
	"auric/internal/providers"
	"auric/internal/providers/consul"
)

var Provider providers.AuricProvider

func InitProvider(providerType string) {
	if providerType == "consul" {
		p := &consul.Consul{}
		err := p.Init()

		if err != nil {
			panic(err)
		}
		Provider = p
	}
}
