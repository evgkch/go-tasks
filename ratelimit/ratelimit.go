package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	request chan chan error
	signal  chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxCount events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	l := &Limiter{
		request: make(chan chan error),
		signal:  make(chan struct{}),
	}
	go l.loop(maxCount, interval)
	return l
}

func (l *Limiter) loop(maxCount int, interval time.Duration) {
	timeline := make([]time.Time, 0, maxCount)
	var timer *time.Timer

	for {
		var request chan chan error
		var timeout <-chan time.Time

		now := time.Now()

		if len(timeline) < maxCount {
			request = l.request
		} else {
			next := timeline[0].Add(interval)
			if now.After(next) {
				request = l.request
			} else {
				if timer == nil {
					timer = time.NewTimer(next.Sub(now))
				}
				timeout = timer.C
			}
		}

		select {
		case <-l.signal:
			if timer != nil {
				timer.Stop()
			}
			for {
				select {
				case ch := <-l.request:
					ch <- ErrStopped
				default:
					return
				}
			}

		case <-timeout:
			timer = nil

		case response := <-request:
			t := time.Now()
			if len(timeline) < maxCount {
				timeline = append(timeline, t)
			} else {
				timeline = append(timeline[1:], t)
			}
			response <- nil
		}
	}
}

func (l *Limiter) Acquire(ctx context.Context) error {
	response := make(chan error, 1)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case l.request <- response:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-response:
		return err
	}
}

func (l *Limiter) Stop() {
	close(l.signal)
}
