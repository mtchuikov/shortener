package postgres

import (
	"context"

	"github.com/mtchuikov/shortener/internal/gen/sqlc/shorturls"
)

func (s *shortURLs) MarkShortURLsAsActive(ctx context.Context, userID string, slugs []string) error {
	_, err := s.queries.BatchUpdateShortURLStatus(ctx,
		shorturls.BatchUpdateShortURLStatusParams{
			Deleted: false,
			UserID:  userID,
			Slugs:   slugs,
		})

	return err
}

func (s *shortURLs) MarkShortURLsAsDeleted(ctx context.Context, userID string, slugs []string) error {
	_, err := s.queries.BatchUpdateShortURLStatus(ctx,
		shorturls.BatchUpdateShortURLStatusParams{
			Deleted: true,
			UserID:  userID,
			Slugs:   slugs,
		})

	return err
}
