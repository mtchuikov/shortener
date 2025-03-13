package inmemory

import "context"

func (r *Repo) GetShortURL(_ context.Context, originalURL string) (string, error) {
	return r.shortURLs[originalURL], nil
}
