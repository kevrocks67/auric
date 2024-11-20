package api

import (
	"auric/internal/api/catalog"
	"auric/internal/api/golden"
	"auric/internal/api/models"
	"auric/internal/providers"
	consulprovider "auric/internal/providers/consul"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var CClient consulprovider.ConsulClient
var GinRouter *gin.Engine

func TestMain(m *testing.M) {
	ctx := context.Background()
	consul, config := initConsul(ctx)
	client, err := consulprovider.NewConsulClient(config)
	if err != nil {
		panic(err)
	}
	CClient.Client = client
	models.InitProvider("consul", &config)
	GinRouter = setupRouter()

	defer consul.Terminate(ctx)

	exitVal := m.Run()
	os.Exit(exitVal)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Catalog API
	router.GET("/catalog", catalog.GetCatalog)
	router.POST("/catalog", catalog.CreateArtifact)
	router.GET("/catalog/:artifact_type/:artifact_guid", catalog.GetArtifact)
	router.PATCH("/catalog/:artifact_type/:artifact_guid", catalog.UpdateArtifact)
	router.DELETE("/catalog/:artifact_type/:artifact_guid", catalog.DeleteArtifact)

	// Golden API
	router.GET("/golden/:artifact_type/:artifact_name/:artifact_channel", golden.GetGoldenArtifact)
	router.PUT("/golden", golden.PromoteGoldenArtifact)
	router.DELETE("/golden/:artifact_type", golden.DeleteGoldenPath)
	router.DELETE("/golden/:artifact_type/:artifact_name", golden.DeleteGoldenPath)
	router.DELETE("/golden/:artifact_type/:artifact_name/:artifact_channel", golden.DeleteGoldenPath)
	return router
}

func startConsulContainer() (testcontainers.Container, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "consul:1.15",
		Name:         "consul-auric-api-test",
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
	w := httptest.NewRecorder()

	testArtifact := models.GoldenArtifact{
		ArtifactUri:        "/catalog/qcow2/rocky9-base/1",
		Channel:            "prod",
		PromotionTimestamp: "2024-06-19T19:14:58Z",
		PromotedBy:         "testUser",
	}

	testArtifactJson, err := json.Marshal(testArtifact)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/golden", strings.NewReader(string(testArtifactJson)))

	if err != nil {
		t.Fatal(err)
	}

	GinRouter.ServeHTTP(w, req)

	var resp models.GoldenArtifact

	err = json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Equal(t, 201, w.Code)
}

func TestGetGoldenArtifact(t *testing.T) {
	// Publish test data
	w := httptest.NewRecorder()

	testArtifact := models.GoldenArtifact{
		ArtifactUri:        "/catalog/qcow2/rocky9-base/1",
		Channel:            "prod",
		PromotionTimestamp: "2024-06-19T19:14:58Z",
		PromotedBy:         "testUser",
	}

	testArtifactJson, err := json.Marshal(testArtifact)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/golden", strings.NewReader(string(testArtifactJson)))

	if err != nil {
		t.Fatal(err)
	}
	GinRouter.ServeHTTP(w, req)

	// Test getting golden artifact
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/golden/qcow2/rocky9-base/prod", strings.NewReader(string(testArtifactJson)))
	if err != nil {
		t.Fatal(err)
	}

	var resp models.GoldenArtifact

	GinRouter.ServeHTTP(w, req)
	respRaw, _ := io.ReadAll(w.Body)
	err = json.Unmarshal(respRaw, &resp)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, testArtifact.ArtifactUri, resp.ArtifactUri)
	assert.Equal(t, testArtifact.Channel, resp.Channel)
	assert.Equal(t, testArtifact.PromotedBy, resp.PromotedBy)
}

func TestDeleteGoldenArtifact(t *testing.T) {
	testCases := []string{
		"/golden/qcow2",
		"/golden/qcow2/rocky9-base",
		"/golden/qcow2/rocky9-base/prod",
	}

	testArtifact := models.GoldenArtifact{
		ArtifactUri:        "/catalog/qcow2/rocky9-base/1",
		Channel:            "prod",
		PromotionTimestamp: "2024-06-19T19:14:58Z",
		PromotedBy:         "testUser",
	}

	testArtifactJson, err := json.Marshal(testArtifact)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			// Publish test data
			w := httptest.NewRecorder()
			req, err := http.NewRequest("PUT", "/golden", strings.NewReader(string(testArtifactJson)))

			if err != nil {
				t.Fatal(err)
			}
			GinRouter.ServeHTTP(w, req)

			// Test delete function
			w = httptest.NewRecorder()
			req, err = http.NewRequest("DELETE", tc, nil)

			if err != nil {
				t.Fatal(err)
			}

			GinRouter.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			require.JSONEq(t, `{"result": true}`, w.Body.String())
		})
	}
}
