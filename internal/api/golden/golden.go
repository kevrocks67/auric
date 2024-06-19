package golden

import (
	"auric/internal/api/constants"
	"auric/internal/api/models"
	"auric/internal/api/util"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetGoldenArtifact(c *gin.Context) {
	var goldenArtifact models.GoldenArtifact
	artifactName := c.Param("artifact_name")
	artifactType := c.Param("artifact_type")
	artifactChannel := c.Param("artifact_channel")

	value, err := models.Provider.Retrieve(fmt.Sprintf("artifacts/golden/%s/%s/%s", artifactType, artifactName, artifactChannel))
	if err != nil {
		panic(err)
	}

	json.Unmarshal(value, &goldenArtifact)

	c.IndentedJSON(http.StatusOK, goldenArtifact)
}
func PromoteGoldenArtifact(c *gin.Context) {
	var newGoldenArtifact models.GoldenArtifact

	if err := c.BindJSON(&newGoldenArtifact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
	}

	t := time.Now().UTC()
	newGoldenArtifact.PromotionTimestamp = t.Format(time.RFC3339)

	artifactType := util.ExtractTypeFromUri(newGoldenArtifact.ArtifactUri)
	artifactName := util.ExtractNameFromUri(newGoldenArtifact.ArtifactUri)

	key := fmt.Sprintf("%s/%s/%s/%s", constants.GoldenArtifactsUri, artifactType, artifactName, newGoldenArtifact.Channel)

	newGoldenArtifactJSON, err := json.Marshal(newGoldenArtifact)
	if err != nil {
		panic(err)
	}

	err = models.Provider.Store(key, newGoldenArtifactJSON)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, newGoldenArtifact)
}
