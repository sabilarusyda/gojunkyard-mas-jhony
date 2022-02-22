package pipeliner

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getdoer(t *testing.T) {
	type x struct{}

	assert.PanicsWithValue(t, `f must be "func([]T) error" or "func(context.Context, []T)" error`, func() { getdoer(1) })
	assert.PanicsWithValue(t, `f must be "func([]T) error" or "func(context.Context, []T)" error`, func() { getdoer(func() {}) })
	assert.PanicsWithValue(t, `f must be "func([]T) error" or "func(context.Context, []T)" error`, func() { getdoer(func() error { return nil }) })

	doer := getdoer(func(ctx context.Context, i []int) error { return nil })
	assert.Nil(t, doer(context.Background(), []*pipelinerCmd{{v: 1, resCh: make(chan error)}, {v: 1, resCh: make(chan error)}}))

	doer = getdoer(func(ctx context.Context, i []int) error { return errors.New("ahuehue") })
	assert.NotNil(t, doer(context.Background(), []*pipelinerCmd{{v: 1, resCh: make(chan error)}, {v: 1, resCh: make(chan error)}}))

	doer = getdoer(func(i []*x) error { return nil })
	assert.Nil(t, doer(context.Background(), []*pipelinerCmd{{v: &x{}, resCh: make(chan error)}}))

	doer = getdoer(func(i []x) error { return nil })
	assert.Nil(t, doer(context.Background(), []*pipelinerCmd{{v: x{}, resCh: make(chan error)}}))

	doer = getdoer(func(i []int) error { return errors.New("ahuehue") })
	assert.NotNil(t, doer(context.Background(), []*pipelinerCmd{{v: 1, resCh: make(chan error)}, {v: 1, resCh: make(chan error)}}))

	doer = getdoer(func(i []int) error { return errors.New("ahuehue") })
	assert.NotNil(t, doer(context.Background(), []*pipelinerCmd{{v: 1, resCh: make(chan error)}}))

	doer = getdoer(func(i []int) error { return nil })
	assert.Panics(t, func() { doer(context.Background(), []*pipelinerCmd{{v: "abc", resCh: make(chan error)}}) })
}

func BenchmarkDoer(b *testing.B) {
	doer := getdoer(func(i []int) error { return nil })
	data := []*pipelinerCmd{{v: 1, resCh: make(chan error)}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doer(context.Background(), data)
	}
}
