package middlewares

import (
	"context"
	"net/http"
)

type ctxKeyError int

const CtxKeyError ctxKeyError = 0

func RequestContextWithError(req *http.Request, err error) {
	ctx := context.WithValue(req.Context(), CtxKeyError, err)
	*req = *req.WithContext(ctx)
}

func errorFromRequestContext(ctx context.Context) error {
	err, ok := ctx.Value(CtxKeyError).(error)
	if ok {
		return err
	}

	return nil
}
