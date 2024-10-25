package catalog

import (
	"auric/internal/api/constants"
	"auric/internal/api/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateArtifact(c *gin.Context) {
	var newArtifact models.Artifact

	// Get request body
	if err := c.BindJSON(&newArtifact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	// Generate server-side parameters
	t := time.Now().UTC()

	artifactUuid, err := uuid.NewV7()
	if err != nil {
		log.Fatal("cannot generate v7 uuid")
	}

	// Populate server-side parameters
	newArtifact.ArtifactGUID = artifactUuid.String()
	newArtifact.ArtifactUri = fmt.Sprintf("%s/%s/%s",
		constants.ArtifactCatalogUri,
		newArtifact.ArtifactType,
		newArtifact.ArtifactGUID,
	)
	newArtifact.CreationTimestamp = t.Format(time.RFC3339)

	// Marshal object and store in backend
	newArtifactJSON, err := json.Marshal(newArtifact)
	if err != nil {
		panic(err)
	}

	err = models.Provider.Store(
		fmt.Sprintf("%s/%s/%s",
			constants.ArtifactCatalogUri,
			newArtifact.ArtifactType,
			newArtifact.ArtifactGUID,
		), newArtifactJSON,
	)
	if err != nil {
		panic(err)
	}

	// Return generated object as JSON
	c.IndentedJSON(http.StatusCreated, newArtifact)
}

func GetCatalog(c *gin.Context) {
	var artifacts []models.Artifact

	pairs, err := models.Provider.List("artifacts/catalog")
	if err != nil {
		panic(err)
	}

	for _, pair := range pairs {
		var artifact models.Artifact
		json.Unmarshal(pair, &artifact)
		artifacts = append(artifacts, artifact)
	}

	c.IndentedJSON(http.StatusOK, artifacts)
}

func GetArtifact(c *gin.Context) {
	var artifact models.Artifact

	artifactType := c.Param("artifact_type")
	artifactGUID := c.Param("artifact_guid")

	value, err := models.Provider.Retrieve(
		fmt.Sprintf("%s/%s/%s",
			constants.ArtifactCatalogUri,
			artifactType,
			artifactGUID,
		),
	)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(value, &artifact)
	c.IndentedJSON(http.StatusOK, artifact)
}

func UpdateArtifact(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, nil)
}

func DeleteArtifact(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, nil)
}
