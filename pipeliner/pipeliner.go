package pipeliner

import (
	"context"
	"sync"
	"time"
)

// Pipeliner ...
type Pipeliner struct {
	window    time.Duration
	limit     int
	doer      func(context.Context, []*pipelinerCmd) error
	spool     sync.Pool
	reqsBufCh chan []*pipelinerCmd
	reqCh     chan *pipelinerCmd
	timeout   time.Duration
}

// Option ...
type Option func(*Pipeliner)

// SetWindow ...
func SetWindow(window time.Duration, limit int) Option {
	return func(pipeliner *Pipeliner) {
		pipeliner.window = window
		pipeliner.limit = limit
	}
}

// SetTimeout ...
func SetTimeout(timeout time.Duration) Option {
	return func(pipeliner *Pipeliner) {
		pipeliner.timeout = timeout
	}
}

// SetConcurrency ...
func SetConcurrency(size int) Option {
	return func(pipeliner *Pipeliner) {
		pipeliner.reqsBufCh = make(chan []*pipelinerCmd, size)
	}
}

// New ...
func New(f interface{}, opts ...Option) *Pipeliner {
	pipeliner := &Pipeliner{doer: getdoer(f), reqCh: make(chan *pipelinerCmd)}
	for _, opt := range opts {
		opt(pipeliner)
	}

	go func() {
		pipeliner.loop()
	}()

	for i := 0; i < cap(pipeliner.reqsBufCh); i++ {
		pipeliner.reqsBufCh <- make([]*pipelinerCmd, 0, pipeliner.limit)
	}

	return pipeliner
}

func (p *Pipeliner) loop() {
	t := time.NewTimer(0)
	t.Stop()

	reqs := <-p.reqsBufCh

	for {
		select {
		case req, ok := <-p.reqCh:
			if !ok {
				continue
			}

			reqs = append(reqs, req)
			if p.limit > 0 && len(reqs) == p.limit {
				t.Stop()
				reqs = p.flush(reqs)
			} else if len(reqs) == 1 {
				t.Reset(p.window)
			}
		case <-t.C:
			reqs = p.flush(reqs)
		}
	}
}

var ifacepool sync.Pool

func (p *Pipeliner) flush(reqs []*pipelinerCmd) []*pipelinerCmd {
	if len(reqs) == 0 {
		return reqs
	}

	go func() {
		defer func() {
			p.reqsBufCh <- reqs[:0]
		}()

		ctx := context.Background()
		if p.timeout > 0 {
			var cancel func()
			ctx, cancel = context.WithTimeout(ctx, p.timeout)
			defer cancel()
		}

		err := p.doer(ctx, reqs)
		for _, req := range reqs {
			req.resCh <- err
		}
	}()

	return <-p.reqsBufCh
}

// Do ...
func (p *Pipeliner) Do(v interface{}) error {
	cmd := getPipelinerCmd()
	cmd.v = v

	p.reqCh <- cmd
	err := <-cmd.resCh

	putPipelinerCmd(cmd)
	return err
}
