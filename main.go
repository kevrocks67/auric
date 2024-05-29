package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Artifact struct {
	ArtifactType       string `json:"artifact_type"`
	ArtifactID         string `json:"artifact_id"`
	PromotionTimestamp string `json:"promotion_timestamp"`
}

func getGoldenArtifact(c *gin.Context) {
	var goldenArtifact Artifact
	artifactType := c.Param("artifact_type")
	artifactChannel := c.Param("artifact_channel")

	pair, err := GetConsulKVPair(fmt.Sprintf("artifacts/golden/%s/%s", artifactType, artifactChannel))
	if err != nil {
		panic(err)
	}

	json.Unmarshal(pair.Value, &goldenArtifact)

	c.IndentedJSON(http.StatusOK, goldenArtifact)
}

func promoteGoldenArtifact(c *gin.Context) {
	var newGoldenArtifact Artifact
	artifactType := c.Param("artifact_type")
	artifactChannel := c.Param("artifact_channel")

	if err := c.BindJSON(&newGoldenArtifact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	t := time.Now().UTC()
	newGoldenArtifact.PromotionTimestamp = t.Format(time.RFC3339)

	newGoldenArtifactJSON, err := json.Marshal(newGoldenArtifact)
	if err != nil {
		panic(err)
	}

	key := fmt.Sprintf("artifacts/golden/%s/%s", artifactType, artifactChannel)

	//TODO make this an interface so that were not tied to consul
	err = CreateConsulKVPair(key, newGoldenArtifactJSON)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newGoldenArtifact)
}

func main() {
	router := gin.Default()
	router.GET("/golden/:artifact_type/:artifact_channel", getGoldenArtifact)
	router.POST("/golden/:artifact_type/:artifact_channel", promoteGoldenArtifact)

	router.Run("0.0.0.0:8080")
}
