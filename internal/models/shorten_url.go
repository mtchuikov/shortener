package models

type ShortenURL string

func NewShortenURL(u string) (ShortenURL, error) {
	return ShortenURL(u), nil
}

func (s ShortenURL) String() string {
	return string(s)
}
