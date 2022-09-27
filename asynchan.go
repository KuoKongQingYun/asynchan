package asynchan

import (
	"context"
	"errors"
)

type Asynchan struct {
	ChanResult chan interface{}
	ChanError  chan error
	Ctx        context.Context
}

func (p *Asynchan) AwaitOne() (interface{}, error) {
	select {
	case result, ok := <-p.ChanResult:
		if ok {
			return result, nil
		} else {
			return nil, errors.New("no result")
		}
	case err := <-p.ChanError:
		return nil, err
	case <-p.Ctx.Done():
		return nil, errors.New("cancelled")
	}
}
func (p *Asynchan) AwaitAll() ([]interface{}, error) {
	results := make([]interface{}, 0)
	for {
		select {
		case result, ok := <-p.ChanResult:
			if ok {
				results = append(results, result)
			} else {
				return results, nil
			}
		case err := <-p.ChanError:
			return results, err
		case <-p.Ctx.Done():
			return results, errors.New("timeout")
		}
	}
}

func New(ctx context.Context) *Asynchan {
	p := &Asynchan{
		ChanResult: make(chan interface{}),
		ChanError:  make(chan error),
		Ctx:        ctx,
	}
	return p
}
