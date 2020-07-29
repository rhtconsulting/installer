package types

type OfflineConfig struct {
	Rhcos struct {
		Architecture string         `json:"architecture"`
		Version string              `json:"version"`
		Provider string             `json:"provider"`
		Url    string               `json:"url"`
	} `json:"rhcos"`

	Ocpmirror struct {
		Ocbin      string           `json:"ocbin"`
		Pullsecret string           `json:"pullsecret"`
		Src        string           `json:"src"`
		Dest       string           `json:"dest"`
	} `json:"ocpmirror"`

	Ocpdistribution struct {
		Destdir string              `json:"destDir"`
		Isofile string              `json:"isoFile"`
	}
}
