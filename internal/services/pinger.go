package services

import "context"

type pinger interface {
	Error() error
}

type pingService struct {
	pinger pinger
}

func NewPinger(pinger pinger) *pingService {
	return &pingService{pinger}
}

func (s *pingService) Serve(ctx context.Context) error {
	return s.pinger.Error()
}
