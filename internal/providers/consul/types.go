package consul

import capi "github.com/hashicorp/consul/api"

type ConsulClient struct {
	Client *capi.Client
}
