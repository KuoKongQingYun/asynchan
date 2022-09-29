package asynchan

import (
	"context"
	"errors"
)

type Asynchan struct {
	ctx        context.Context
	chanResult chan interface{}
	chanError  chan error
}

func (p *Asynchan) AwaitOne() (interface{}, error) {
	select {
	case result, ok := <-p.chanResult:
		if ok {
			return result, nil
		} else {
			return nil, errors.New("no result")
		}
	case err := <-p.chanError:
		return nil, err
	case <-p.ctx.Done():
		return nil, errors.New("cancelled")
	}
}
func (p *Asynchan) AwaitAll() ([]interface{}, error) {
	results := make([]interface{}, 0)
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
func (p *Asynchan) SetResult(result interface{}) error {
	select {
	case p.chanResult <- result:
		return nil
	case <-p.ctx.Done():
		return errors.New("cancelled")
	}
}
func (p *Asynchan) SetError(err error) error {
	select {
	case p.chanError <- err:
		return nil
	case <-p.ctx.Done():
		return errors.New("cancelled")
	}
}
func (p *Asynchan) Close(err error) {
	close(p.chanResult)
	close(p.chanError)
}
func New(ctx context.Context) *Asynchan {
	p := &Asynchan{
		chanResult: make(chan interface{}),
		chanError:  make(chan error),
		ctx:        ctx,
	}
	return p
}
