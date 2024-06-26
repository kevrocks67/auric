package api

import (
	"context"
	"encoding/json"
	"fmt"
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

func TestMain(m *testing.M) {
	ctx := context.Background()
	consul, baseUrl := initConsul(ctx)
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
		Image:        "consul:latest",
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

func initConsul(ctx context.Context) (testcontainers.Container, string) {
	consul, err := startConsulContainer()
	if err != nil {
		panic("Could not start Consul container: " + err.Error())
	}

	ip, err := consul.Host(ctx)
	if err != nil {
		panic("Could not get container host" + err.Error())
	}

	port, err := consul.MappedPort(ctx, "9200")
	if err != nil {
		panic("Could not retrive the mapped port: " + err.Error())
	}

	baseUrl := fmt.Sprintf("http://%s:%s", ip, port.Port())

	return consul, baseUrl
}

func TestPromoteGoldenArtifact(t *testing.T) {
	router := setupRouter()
	router.POST("/artifacts/golden/:artifact_type/:artifact_channel", promoteGoldenArtifact)

	w := httptest.NewRecorder()

	testArtifact := Artifact{
		ArtifactType: "test",
		ArtifactID:   "1234",
	}

	testArtifactJson, err := json.Marshal(testArtifact)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/artifacts/golden/test/prod", strings.NewReader(string(testArtifactJson)))
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
