package cli

type Album struct {
	AlbumType            string         `json:"album_type"`
	TotalTracks          int64          `json:"total_tracks"`
	ExternalURLs         ExternalURL    `json:"external_urls"`
	Href                 string         `json:"href"`
	ID                   string         `json:"id"`
	Name                 string         `json:"name"`
	ReleaseDate          string         `json:"release_date"`
	ReleaseDatePrecision string         `json:"release_date_precision"`
	Restrictions         Restriction    `json:"restrictions"`
	Type                 string         `json:"type"`
	URI                  string         `json:"uri"`
	Artists              []Artist       `json:"artists"`
	Tracks               []Track        `json:"tracks"`
	Copyrights           []Copyright    `json:"copyrights"`
	ExternalIDs          ExternalIDList `json:"external_ids"`
}

type ExternalURL struct {
	Spotify string `json:"spotify"`
}

type Restriction struct {
	Reason string `json:"reason"`
}

type ExternalIDList struct {
	ISRC string `json:"isrc"`
	EAN  string `json:"ean"`
	UPC  string `json:"upc"`
}

type Copyright struct {
	Text string `json:"text"`
	Type string `json:"type"`
}
