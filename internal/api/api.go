package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	ac "auric/internal/consul"

	"github.com/gin-gonic/gin"
)

func GetGoldenArtifact(c *gin.Context) {
	var goldenArtifact GoldenArtifact
	artifactName := c.Param("artifact_name")
	artifactType := c.Param("artifact_type")
	artifactChannel := c.Param("artifact_channel")

	pair, err := ac.GetConsulKVPair(fmt.Sprintf("artifacts/golden/%s/%s/%s", artifactType, artifactName, artifactChannel))
	if err != nil {
		panic(err)
	}

	json.Unmarshal(pair.Value, &goldenArtifact)

	c.IndentedJSON(http.StatusOK, goldenArtifact)
}

func PromoteGoldenArtifact(c *gin.Context) {
	var newGoldenArtifact GoldenArtifact

	if err := c.BindJSON(&newGoldenArtifact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	t := time.Now().UTC()
	newGoldenArtifact.PromotionTimestamp = t.Format(time.RFC3339)

	artifactType := extractTypeFromUri(newGoldenArtifact.ArtifactUri)
	artifactName := extractNameFromUri(newGoldenArtifact.ArtifactUri)

	key := fmt.Sprintf("%s/%s/%s/%s", goldenArtifactsUri, artifactType, artifactName, newGoldenArtifact.Channel)

	newGoldenArtifactJSON, err := json.Marshal(newGoldenArtifact)
	if err != nil {
		panic(err)
	}

	//TODO make this an interface so that were not tied to consul
	err = ac.CreateConsulKVPair(key, newGoldenArtifactJSON)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newGoldenArtifact)
}

func extractNameFromUri(uri string) string {
	return strings.Split(uri, "/")[3]
}

func extractTypeFromUri(uri string) string {
	return strings.Split(uri, "/")[2]
}

func CreateArtifact(c *gin.Context) {
	var newArtifact Artifact

	if err := c.BindJSON(&newArtifact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}
	t := time.Now().UTC()

	newArtifact.ArtifactUri = fmt.Sprintf("%s/%s/%s/%s", artifactCatalogUri, newArtifact.ArtifactType, newArtifact.ArtifactName, newArtifact.ArtifactId)
	newArtifact.CreationTimestamp = t.Format(time.RFC3339)

	newArtifactJSON, err := json.Marshal(newArtifact)
	if err != nil {
		panic(err)
	}

	//TODO make this an interface so that were not tied to consul
	err = ac.CreateConsulKVPair(newArtifact.ArtifactUri, newArtifactJSON)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newArtifact)

}

func Serve(args ...string) {
	addr := "0.0.0.0:8080"

	if len(args) > 1 {
		addr = args[1]
	}

	router := gin.Default()
	router.GET("/golden/:artifact_type/:artifact_name/:artifact_channel", GetGoldenArtifact)
	router.POST("/golden", PromoteGoldenArtifact)
	router.POST("/catalog", CreateArtifact)

	router.Run(addr)
}
