package models

type ShortURLRequest struct {
	URL string `json:"url"`
}

type ShortURLResponse struct {
	Result string `json:"result"`
}

type BatchShortURLs []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchShortenURLs []struct {
	CorrelationID string `json:"correlation_id"`
	ShortenURL    string `json:"short_url"`
}
