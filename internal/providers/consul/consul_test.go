package consul

import (
	"auric/internal/providers"
	"context"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
)

var Client ConsulClient

func TestMain(m *testing.M) {
	ctx := context.Background()
	consul, config := initConsul(ctx)
	client, err := newConsulClient(config)
	if err != nil {
		panic(err)
	}
	Client.client = client

	defer consul.Terminate(ctx)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func startConsulContainer() (testcontainers.Container, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "consul:1.15",
		Name:         "consul-auric-test",
		ExposedPorts: []string{"8500/tcp"},
		Cmd:          []string{"agent", "-dev", "-client", "0.0.0.0"},
		WaitingFor:   wait.NewHTTPStrategy("/v1/status/leader"),
	}

	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)

	return container, err
}

func initConsul(ctx context.Context) (testcontainers.Container, providers.ProviderConfig) {
	consul, err := startConsulContainer()
	if err != nil {
		panic("Could not start Consul container: " + err.Error())
	}

	ip, err := consul.Host(ctx)
	if err != nil {
		panic("Could not get container host" + err.Error())
	}

	port, err := consul.MappedPort(ctx, "8500")
	if err != nil {
		panic("Could not retrive the mapped port: " + err.Error())
	}

	config := providers.ProviderConfig{
		ConsulConfig: struct {
			Address    string `json:"address"`
			Port       string `json:"port"`
			Datacenter string `json:"datacenter"`
		}{
			Address:    ip,
			Port:       port.Port(),
			Datacenter: "dc1",
		},
	}

	return consul, config
}

func TestCreateConsulKVPair(t *testing.T) {
	err := CreateConsulKVPair(&Client, "testKey", []byte("something"))
	assert.Equal(t, err, nil)
}

func TestGetConsulKVPair(t *testing.T) {
	pair, err := GetConsulKVPair(&Client, "testKey")
	assert.Equal(t, string(pair.Value), "something")
	assert.Equal(t, err, nil)
}
