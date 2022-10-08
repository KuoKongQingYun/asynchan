package asynchan

import (
	"context"
	"errors"
)

type Asynchan[T any] struct {
	ctx        context.Context
	chanResult chan T
	chanError  chan error
}

func (p *Asynchan[T]) AwaitOne() (T, error) {
	var result T
	select {
	case result, ok := <-p.chanResult:
		if ok {
			return result, nil
		} else {
			return result, errors.New("no result")
		}
	case err := <-p.chanError:
		return result, err
	case <-p.ctx.Done():
		return result, errors.New("cancelled")
	}
}
func (p *Asynchan[T]) AwaitAll() ([]T, error) {
	results := make([]T, 0)
	for {
		select {
		case result, ok := <-p.chanResult:
			if ok {
				results = append(results, result)
			} else {
				return results, nil
			}
		case err := <-p.chanError:
			return results, err
		case <-p.ctx.Done():
			return results, errors.New("timeout")
		}
	}
}
func (p *Asynchan[T]) SetResult(result T) error {
	select {
	case p.chanResult <- result:
		return nil
	case <-p.ctx.Done():
		return errors.New("cancelled")
	}
}
func (p *Asynchan[T]) SetError(err error) error {
	select {
	case p.chanError <- err:
		return nil
	case <-p.ctx.Done():
		return errors.New("cancelled")
	}
}
func (p *Asynchan[T]) Close() {
	close(p.chanResult)
	close(p.chanError)
}
func New[T any](ctx context.Context) *Asynchan[T] {
	p := &Asynchan[T]{
		chanResult: make(chan T),
		chanError:  make(chan error),
		ctx:        ctx,
	}
	return p
}
