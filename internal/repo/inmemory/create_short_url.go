package inmemory

import "context"

func (r *Repo) CreateShortURL(_ context.Context, originalURL, shortURL string) error {
	r.mu.Lock()
	r.originalURLs[shortURL] = originalURL
	r.shortURLs[originalURL] = shortURL
	r.mu.Unlock()
	return nil
}
