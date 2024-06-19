package catalog

import (
	"auric/internal/api/constants"
	"auric/internal/api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateArtifact(c *gin.Context) {
	var newArtifact models.Artifact

	if err := c.BindJSON(&newArtifact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}
	t := time.Now().UTC()

	newArtifact.ArtifactUri = fmt.Sprintf("%s/%s/%s/%s", constants.ArtifactCatalogUri, newArtifact.ArtifactType, newArtifact.ArtifactName, newArtifact.ArtifactId)
	newArtifact.CreationTimestamp = t.Format(time.RFC3339)

	newArtifactJSON, err := json.Marshal(newArtifact)
	if err != nil {
		panic(err)
	}

	err = models.Provider.Store(newArtifact.ArtifactUri, newArtifactJSON)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newArtifact)
}
