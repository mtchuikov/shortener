package models

type ShortURL struct {
	URL string `json:"url"`
}

type ShortURLResult struct {
	Result string `json:"result"`
}

type BatchCreateShortURLs struct {
	Slug        string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type BatchCreateShortURLsResult struct {
	Slug     string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type ListShortURLsByUserResult struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
