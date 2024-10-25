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

/* GetGoldenArtifact returns a GoldenArtifact object
*  It gets the data based on parameters in the URI
*  An example URI for an artifact named rocky9-base
*  of type qcow2 in prod would be obtained by making
*  a GET request to /golden/qcow2/rocky9-base/prod
 */
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

	artifactType := util.ExtractTypeFromGoldenUri(newGoldenArtifact.ArtifactUri)
	artifactName := util.ExtractNameFromGoldenUri(newGoldenArtifact.ArtifactUri)

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

func DeleteGoldenPath(c *gin.Context) {
	c.IndentedJSON(http.StatusNotImplemented, nil)
}
