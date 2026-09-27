package cli

type Track struct {
	Album        Album          `json:"album"`
	Artists      []Artist       `json:"artists"`
	DiscNumber   int64          `json:"disc_number"`
	DurationMs   int64          `json:"duration_ms"`
	Explicit     bool           `json:"explicit"`
	ExternalIDs  ExternalIDList `json:"external_ids"`
	ExternalURLs ExternalURL    `json:"external_url"`
	Href         string         `json:"href"`
	ID           string         `json:"id"`
	IsPlayable   bool           `json:"is_playable"`
	Restrictions Restriction    `json:"restrictions"`
	Name         string         `json:"name"`
	TrackNumber  int64          `json:"track_number"`
	Type         string         `json:"type"`
	URI          string         `json:"uri"`
	IsLocal      bool           `json:"is_local"`
}
