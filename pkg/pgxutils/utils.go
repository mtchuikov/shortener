package pgxutils

import (
	"context"
	"fmt"
	"time"
)

const pingTimeout = 2 * time.Second

func pingWithTimeout(
	ctx context.Context,
	ping interface {
		Ping(ctx context.Context) error
	},
) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	err := ping.Ping(timeoutCtx)
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrFailedToPing, err)
		return err
	}

	return nil
}
