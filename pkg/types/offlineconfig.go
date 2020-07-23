package types

type OfflineConfig struct {
	Rhcos struct {
		Architecture string         `json:"architecture"`
		Version string              `json:"version"`
		Provider string             `json:"provider"`
		Url    string               `json:"url"`
	} `json:"rhcos"`
}
