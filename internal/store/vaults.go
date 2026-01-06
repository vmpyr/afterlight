package store

import "time"

type ListArtifactsResponse struct {
	VaultName string     `json:"vault_name"`
	Hint      string     `json:"hint,omitempty"`
	Artifacts []Artifact `json:"artifacts"`
	CreatedAt time.Time  `json:"created_at"`
}
