package timeout

import (
	"context"
	"time"

	"github.com/hokonco/gdk/errors"
)

// WithContext ...
func WithContext(ctx context.Context, fns ...func() error) error {
	if len(fns) < 1 {
		return errors.New("timeout: empty func")
	}
	var waitingFor = 0
	var c = make(chan error, 1)
	for i := range fns {
		if fns[i] == nil {
			continue
		}
		waitingFor++
		go func(i int) { c <- fns[i]() }(i)
	}
	if waitingFor < 1 {
		return errors.New("timeout: empty func")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}

// WithDuration ...
func WithDuration(d time.Duration, fn func() error) error {
	var ctx, cancel = context.WithTimeout(context.Background(), d)
	defer cancel()
	return WithContext(ctx, fn)
}

// WithDurationString ...
func WithDurationString(s string, fn func() error) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	return WithDuration(d, fn)
}
