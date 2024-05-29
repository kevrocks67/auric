package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	capi "github.com/hashicorp/consul/api"
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

	client, err := capi.NewClient(capi.DefaultConfig())
	if err != nil {
		panic(err)
	}

	kv := client.KV()

	pair, _, err := kv.Get(fmt.Sprintf("artifacts/golden/%s/%s", artifactType, artifactChannel), nil)
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

	client, err := capi.NewClient(capi.DefaultConfig())
	if err != nil {
		panic(err)
	}
	kv := client.KV()

	newGoldenArtifactJSON, err := json.Marshal(newGoldenArtifact)
	if err != nil {
		panic(err)
	}

	p := &capi.KVPair{
		Key:   fmt.Sprintf("artifacts/golden/%s/%s", artifactType, artifactChannel),
		Value: newGoldenArtifactJSON,
	}
	_, err = kv.Put(p, nil)
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
