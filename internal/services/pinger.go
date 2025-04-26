package services

import "context"

type ping interface {
	Error() error
}

var _ Pinger = (*pinger)(nil)

type pinger struct {
	ping ping
}

func NewPinger(ping ping) *pinger {
	return &pinger{ping}
}

func (s *pinger) Ping(ctx context.Context) error {
	return s.ping.Error()
}
