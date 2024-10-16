package models

type Artifact struct {
	ArtifactGUID      string   `json:"artifact_guid"`
	ArtifactUri       string   `json:"artifact_uri"`
	ArtifactType      string   `json:"artifact_type"`
	ArtifactName      string   `json:"artifact_name"`
	ArtifactId        string   `json:"artifact_id"`
	CreationTimestamp string   `json:"creation_timestamp"`
	UploadedBy        string   `json:"uploaded_by"`
	HasChildren       bool     `json:"has_children"`
	Parent            string   `json:"parent"`
	Children          []string `json:"children"`
}

type GoldenArtifact struct {
	ArtifactUri        string `json:"artifact_uri"`
	Channel            string `json:"channel"`
	PromotionTimestamp string `json:"promotion_timestamp"`
	PromotedBy         string `json:"promoted_by"`
}
