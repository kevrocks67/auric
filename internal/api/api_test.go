package api

import (
	"auric/internal/api/golden"
	"auric/internal/api/models"
	"auric/internal/providers"
	consulprovider "auric/internal/providers/consul"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var CClient consulprovider.ConsulClient

func TestMain(m *testing.M) {
	ctx := context.Background()
	consul, config := initConsul(ctx)
	client, err := consulprovider.NewConsulClient(config)
	if err != nil {
		panic(err)
	}
	CClient.Client = client

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

func TestPromoteGoldenArtifact(t *testing.T) {
	router := setupRouter()
	router.POST("/artifacts/golden/:artifact_type/:artifact_channel", golden.PromoteGoldenArtifact)

	w := httptest.NewRecorder()

	testArtifact := models.GoldenArtifact{
		ArtifactUri:        "catalog/qcow2/rocky9-base/1",
		Channel:            "prod",
		PromotionTimestamp: "2024-06-19T19:14:58Z",
		PromotedBy:         "testUser",
	}

	testArtifactJson, err := json.Marshal(testArtifact)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/golden/qcow2/rocky9-base/prod", strings.NewReader(string(testArtifactJson)))
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	var resp models.GoldenArtifact

	err = json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, testArtifact.ArtifactUri, resp.ArtifactUri)
	assert.Equal(t, testArtifact.Channel, resp.Channel)
	assert.Equal(t, testArtifact.PromotedBy, resp.PromotedBy)

	assert.Equal(t, 201, w.Code)
}
