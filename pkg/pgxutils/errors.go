package pgxutils

import "errors"

var ErrFailedToPing = errors.New("failed to ping postgres")
