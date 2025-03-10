package inmemory

import "context"

func (r *Repo) GetOriginalURL(_ context.Context, shortURL string) (string, error) {
	return r.originalURLs[shortURL], nil
}
