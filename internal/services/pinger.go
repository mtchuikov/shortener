package services

import (
	"context"
)

type ping interface {
	Error() error
}

type pinger struct{ ping ping }

func NewPinger(ping ping) *pinger {
	return &pinger{ping: ping}
}

func (s *pinger) Serve(ctx context.Context) error {
	return s.ping.Error()
}
