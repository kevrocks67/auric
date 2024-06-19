package api

import (
	"auric/internal/api/catalog"
	"auric/internal/api/golden"
	"auric/internal/api/models"

	"github.com/gin-gonic/gin"
)

func Serve(args ...string) {
	addr := "0.0.0.0:8080"

	if len(args) > 1 {
		addr = args[1]
	}

	models.InitProvider("consul")

	router := gin.Default()
	router.GET("/golden/:artifact_type/:artifact_name/:artifact_channel", golden.GetGoldenArtifact)
	router.POST("/golden", golden.PromoteGoldenArtifact)
	router.POST("/catalog", catalog.CreateArtifact)

	router.Run(addr)
}
