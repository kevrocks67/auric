package consul

import capi "github.com/hashicorp/consul/api"

type Consul struct {
	client *capi.Client
}
