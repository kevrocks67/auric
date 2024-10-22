package api

import (
	"auric/internal/api/catalog"
	"auric/internal/api/golden"
	"auric/internal/api/models"
	"auric/internal/providers"

	"github.com/gin-gonic/gin"
)

func Serve(args ...string) {
	addr := "0.0.0.0:8080"

	if len(args) > 1 {
		addr = args[1]
		//TODO implement config file getting loaded and determine provider
		// configFile = args[2]
	}

	config := providers.ProviderConfig{
		ConsulConfig: struct {
			Address    string `json:"address"`
			Port       string `json:"port"`
			Datacenter string `json:"datacenter"`
		}{
			Address:    "127.0.0.1",
			Port:       "8500",
			Datacenter: "dc1",
		},
	}

	models.InitProvider("consul", &config)

	router := gin.Default()
	router.GET("/golden/:artifact_type/:artifact_name/:artifact_channel", golden.GetGoldenArtifact)
	router.POST("/golden", golden.PromoteGoldenArtifact)
	router.GET("/catalog", catalog.GetCatalog)
	router.POST("/catalog", catalog.CreateArtifact)

	router.Run(addr)
}
