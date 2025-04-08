package models

type URLToShort struct {
	URL string `json:"url"`
}

type ShortenURL struct {
	Result string `json:"result"`
}
