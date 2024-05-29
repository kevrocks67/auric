package consul

import capi "github.com/hashicorp/consul/api"

type ConsulClient struct {
	client *capi.Client
}
