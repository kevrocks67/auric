package api

import (
	"auric/internal/api/golden"
	"auric/internal/api/models"
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
	router.POST("/artifacts/golden/:artifact_type/:artifact_channel", golden.PromoteGoldenArtifact)

	w := httptest.NewRecorder()

	testArtifact := models.Artifact{
		ArtifactType: "qcow2",
		ArtifactName: "rocky9-base",
		ArtifactId:   "1",
	}

	testArtifactJson, err := json.Marshal(testArtifact)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/artifacts/golden", strings.NewReader(string(testArtifactJson)))
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestGetGoldenArtifact(t *testing.T) {
	router := setupRouter()
	router.GET("/artifacts/golden/qcow2/rocky9-base/prod", golden.GetGoldenArtifact)

	w := httptest.NewRecorder()

	testArtifact := models.GoldenArtifact{
		ArtifactUri:        "artifacts/catalog/qcow2/rocky9-base/1",
		Channel:            "prod",
		PromotionTimestamp: "2024-06-19T19:14:58Z",
		PromotedBy:         "testUser",
	}

	testArtifactJson, err := json.Marshal(testArtifact)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/artifacts/golden/qcow2/rocky9-base/prod", strings.NewReader(string(testArtifactJson)))
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}
