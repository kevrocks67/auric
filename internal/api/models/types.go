package models

type Artifact struct {
	ArtifactGUID      string   `json:"artifact_guid"`
	ArtifactUri       string   `json:"artifact_uri"`
	ArtifactType      string   `json:"artifact_type"`
	ArtifactName      string   `json:"artifact_name"`
	ArtifactVersion   string   `json:"artifact_version"`
	CreationTimestamp string   `json:"creation_timestamp"`
	LastUpdated       string   `json:"last_updated"`
	UploadedBy        string   `json:"uploaded_by"`
	HasChildren       bool     `json:"has_children"`
	Parent            string   `json:"parent"`
	Children          []string `json:"children"`
}

type ArtifactPatch struct {
	ArtifactType    string `json:"artifact_type"`
	ArtifactName    string `json:"artifact_name"`
	ArtifactVersion string `json:"artifact_version"`
}

type GoldenArtifact struct {
	ArtifactUri        string `json:"artifact_uri"`
	Channel            string `json:"channel"`
	PromotionTimestamp string `json:"promotion_timestamp"`
	PromotedBy         string `json:"promoted_by"`
}
