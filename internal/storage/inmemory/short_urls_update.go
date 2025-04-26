package inmemory

import (
	"context"
	"fmt"

	"github.com/mtchuikov/shortener/internal/storage"
)

func (s *shortURLs) MarkShortURLsAsActive(ctx context.Context, userID string, slugs []string) error {
	err := fmt.Errorf("%s: unimplemented", storage.ErrMsgSomethingWentWrong)
	return err
}

func (s *shortURLs) MarkShortURLsAsDeleted(ctx context.Context, userID string, slugs []string) error {
	err := fmt.Errorf("%s: unimplemented", storage.ErrMsgSomethingWentWrong)
	return err
}
