package pipeliner

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.Panics(t, func() { New(func() {}) })

	pipe := New(func([]int) error { return errors.New("TESTING_NEW") }, SetConcurrency(2), SetTimeout(time.Second), SetWindow(100*time.Microsecond, 2))
	assert.Equal(t, time.Microsecond*100, pipe.window)
	assert.Equal(t, time.Second, pipe.timeout)
	assert.Equal(t, 2, cap(pipe.reqsBufCh))
	assert.Equal(t, 2, pipe.limit)
	assert.Equal(t, errors.New("TESTING_NEW"), pipe.doer(context.Background(), nil))
}

func Test_Limit(t *testing.T) {
	pipe := New(func([]int) error { return errors.New("AHUEHUE") }, SetConcurrency(2), SetTimeout(time.Second), SetWindow(time.Millisecond, 2))
	now := time.Now()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		assert.Equal(t, errors.New("AHUEHUE"), pipe.Do(2))
	}()
	go func() {
		defer wg.Done()
		assert.Equal(t, errors.New("AHUEHUE"), pipe.Do(1))
	}()
	wg.Wait()

	assert.False(t, time.Since(now) > time.Second, "Queue must be less than 1 second")
}

func Test_Window(t *testing.T) {
	pipe := New(func([]int) error { return errors.New("AHUEHUE") }, SetConcurrency(2), SetTimeout(time.Second), SetWindow(100*time.Microsecond, 2))

	now := time.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		assert.Equal(t, errors.New("AHUEHUE"), pipe.Do(2))
	}()
	wg.Wait()

	assert.True(t, time.Since(now) > 100*time.Microsecond, "Queue must be more than 1 microsecond")
}
