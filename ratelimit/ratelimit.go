package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	tokens chan struct{}
	stamps chan struct{}
	signal chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxCount events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	l := &Limiter{
		tokens: make(chan struct{}, maxCount),
		stamps: make(chan struct{}, maxCount),
		signal: make(chan struct{}),
	}
	for i := 0; i < maxCount; i++ {
		l.tokens <- struct{}{}
	}
	go l.loop(interval)
	return l
}

func (l *Limiter) loop(interval time.Duration) {
	var scheduler []*time.Timer
	var timeout <-chan time.Time

	for {
		if len(scheduler) > 0 {
			timeout = scheduler[0].C
		}

		select {
		case <-l.stamps:
			scheduler = append(scheduler, time.NewTimer(interval))
		case <-timeout:
			scheduler = scheduler[1:]
			l.tokens <- struct{}{}
		case <-l.signal:
			for _, timer := range scheduler {
				timer.Stop()
			}
			return
		}
	}
}

func (l *Limiter) Acquire(ctx context.Context) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-l.tokens:
		l.stamps <- struct{}{}
		select {
		case <-l.signal:
			return ErrStopped
		default:
			return nil
		}
	case <-l.signal:
		return ErrStopped
	}
}

func (l *Limiter) Stop() {
	close(l.signal)
}
