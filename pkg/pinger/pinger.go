package pinger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mtchuikov/shortener/pkg/backoff"
	"github.com/mtchuikov/shortener/pkg/chsubscription"
)

const (
	DefaultPingInterval    = 3 * time.Second
	DefaultBackoffInterval = 2 * time.Second
	DefaultBackoffMax      = 6 * time.Second
	DefaultFailedThreshold = 0
)

type pinger interface {
	Ping(ctx context.Context) error
}

type Pinger struct {
	mu              sync.Mutex
	pinger          pinger
	ticker          *time.Ticker
	backoff         *backoff.Backoff
	err             error
	failedThreshold uint8
	failedCount     uint8
	chsub           *chsubscription.ChSubscription[error]
	closeOnce       sync.Once
}

func New(pinger pinger, opts ...Option) *Pinger {
	p := &Pinger{
		mu:              sync.Mutex{},
		pinger:          pinger,
		ticker:          time.NewTicker(DefaultPingInterval),
		err:             nil,
		chsub:           chsubscription.New[error](),
		closeOnce:       sync.Once{},
		failedThreshold: DefaultFailedThreshold,
	}

	backoff := backoff.New()
	backoff.Interval = DefaultBackoffInterval
	backoff.Max = DefaultBackoffMax

	p.backoff = backoff
	p.failedCount = 0

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Pinger) Ping(ctx context.Context, timeout time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				p.chsub.Notify(ctx, ctx.Err())
				return
			case <-p.ticker.C:

				timeoutCtx, cancel := context.WithTimeout(ctx, timeout)

				p.mu.Lock()
				err := p.pinger.Ping(timeoutCtx)
				p.mu.Unlock()

				cancel()

				if err != nil {
					if p.failedThreshold > 0 &&
						p.failedCount >= p.failedThreshold {
						err = fmt.Errorf("%w: %w", ErrFailedThresholdExceeded, err)
					}

					p.err = err
					p.chsub.Notify(ctx, err)

					delay := p.backoff.Next()
					select {
					case <-ctx.Done():
						return
					case <-time.After(delay):
						continue
					}
				}

				p.backoff.Reset()
				p.failedCount = 0
			}
		}
	}()
}

func (p *Pinger) Replace(pinger pinger) {
	p.mu.Lock()
	p.pinger = pinger
	p.mu.Unlock()
}

func (p *Pinger) Error() error {
	return p.err
}

func (p *Pinger) Subscribe() <-chan error {
	return p.chsub.Subscribe(3)
}

func (p *Pinger) Unsubscribe(item <-chan error) {
	p.chsub.Unsubscribe(item)
}

func (p *Pinger) Close(ctx context.Context) error {
	p.closeOnce.Do(
		func() {
			p.ticker.Stop()
			p.chsub.Close()
		})

	return nil
}
