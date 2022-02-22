package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemory(t *testing.T) {
	im := NewInMemory(SetInMemoryCapacity(10))
	assert.Equal(t, &InMemory{
		storage: make(map[uint64]*inMemoryObject, 10),
	}, im)
}

func TestSetInMemoryCapacity(t *testing.T) {
	im := NewInMemory()
	im.Set(1234, &object{}, time.Minute)

	assert.NotNil(t, im)
	assert.Len(t, im.storage, 1)

	SetInMemoryCapacity(10)(im)

	assert.NotNil(t, im)
	assert.Len(t, im.storage, 1)
}

func TestInMemory_GetSetDelete(t *testing.T) {
	const key = uint64(123)
	obj := &object{}
	im := NewInMemory()

	// condition 1. not exist on memory
	r, e := im.Get(key)
	assert.Nil(t, e)
	assert.Nil(t, r)

	// condition 2. the cache has been expire
	now = func() time.Time { return time.Date(2010, time.February, 1, 0, 0, 0, 0, time.UTC) }
	im.Set(key, obj, time.Minute)
	now = func() time.Time { return time.Date(2011, time.February, 1, 0, 0, 0, 0, time.UTC) }
	r, e = im.Get(key)
	assert.Nil(t, e)
	assert.Nil(t, r)

	// condition 3. valid cache
	im.Set(key, obj, time.Minute)
	r, e = im.Get(key)
	assert.Nil(t, e)
	assert.Equal(t, obj, r)

	// condition 4. cache deleted
	im.Delete(key)
	r, e = im.Get(key)
	assert.Nil(t, e)
	assert.Nil(t, r)

	// condition 5. set expire < 0
	im.Set(key, obj, time.Duration(-100))
	r, e = im.Get(key)
	assert.Nil(t, e)
	assert.Nil(t, r)

	// condition 6. cache without expire
	im.Set(key, obj, time.Duration(0))
	r, e = im.Get(key)
	assert.Nil(t, e)
	assert.Equal(t, obj, r)
}

// BenchmarkInMemory_Get-4   	50000000	        25.7 ns/op	       0 B/op	       0 allocs/op
func BenchmarkInMemory_Get(b *testing.B) {
	const key = 123456789
	im := NewInMemory()
	im.Set(key, &object{}, 0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		im.Get(123)
	}
}

// BenchmarkInMemory_Set-4   	20000000	        63.9 ns/op	      16 B/op	       1 allocs/op
func BenchmarkInMemory_Set(b *testing.B) {
	const key = 123456789
	im := NewInMemory()
	obj := object{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		im.Set(key, &obj, 0)
	}
}

// BenchmarkInMemory_Delete-4   	50000000	        31.4 ns/op	       0 B/op	       0 allocs/op
func BenchmarkInMemory_Delete(b *testing.B) {
	const key = 123456789
	im := NewInMemory()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		im.Delete(key)
	}
}
