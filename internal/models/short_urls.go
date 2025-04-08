package models

type URLsToShort []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortID       string `json:"-"`
}

type ShortenURLs []struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"original_url"`
}
