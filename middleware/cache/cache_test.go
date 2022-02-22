package cache

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_generateKey(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "https://supersoccer.tv/api/v2/campaigns/banner-positions/web-main-landing-page", nil)
	key := generateKey(request)
	assert.Equal(t, uint64(1149593713296141760), key)

	request = httptest.NewRequest(http.MethodGet, "https://supersoccer.tv/api/v2/campaigns/banner-positions/web-main-landing-page?include=banners", nil)
	key = generateKey(request)
	assert.Equal(t, uint64(14995860473677391509), key)
}

// Benchmark_generateKey-4   	 3000000	       450 ns/op	     280 B/op	       4 allocs/op
// allocation: r.URL.String()
func Benchmark_generateKey(b *testing.B) {
	request := httptest.NewRequest(http.MethodGet, "https://supersoccer.tv/api/v2/campaigns/banner-positions/web-main-landing-page?include=banners", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateKey(request)
	}
}

func Test_normalize(t *testing.T) {
	requestOne := httptest.NewRequest(http.MethodGet, "https://supersoccer.tv/api/v2/campaigns/banner-positions/web-main-landing-page?include=banners&data[]=1&data[]=2", nil)
	normalize(requestOne.URL)

	requestTwo := httptest.NewRequest(http.MethodGet, "https://supersoccer.tv/api/v2/campaigns/banner-positions/web-main-landing-page?data[]=2&data[]=1&include=banners", nil)
	normalize(requestTwo.URL)

	assert.Equal(t, requestOne.URL, requestTwo.URL)
}

// Benchmark_normalize-4   	 1000000	      1439 ns/op	     784 B/op	      19 allocs/op
func Benchmark_normalize(b *testing.B) {
	request := httptest.NewRequest(http.MethodGet, "https://supersoccer.tv/api/v2/campaigns/banner-positions/web-main-landing-page?include=banners&data[]=1&data[]=2", nil)
	for i := 0; i < b.N; i++ {
		normalize(request.URL)
	}
}

func Test_renderResponse(t *testing.T) {
	var (
		code   = 200
		body   = []byte("{}")
		header = http.Header(map[string][]string{
			"Content-Type": []string{"application/json"},
		})
	)
	recorder := httptest.NewRecorder()
	obj := object{
		Header:     header,
		StatusCode: code,
		Body:       body,
	}
	renderResponse(recorder, &obj)

	assert.Equal(t, code, recorder.Code)
	assert.Equal(t, body, recorder.Body.Bytes())
	assert.Equal(t, header, recorder.Header())
}

func Benchmark_renderResponse(b *testing.B) {
	var (
		code   = 200
		body   = []byte("{}")
		header = http.Header(map[string][]string{
			"Content-Type": []string{"application/json"},
		})
	)
	obj := object{
		Header:     header,
		StatusCode: code,
		Body:       body,
	}
	recorder := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderResponse(recorder, &obj)
	}
}

func Test_cache_SetStorage(t *testing.T) {
	storage := NewInMemory()

	c := new(cache)
	assert.Equal(t, &cache{}, c)

	c.SetStorage(storage)
	assert.Equal(t, &cache{storage: storage}, c)
}

func Test_cache_SetTTL(t *testing.T) {
	c := new(cache)
	assert.Equal(t, &cache{}, c)

	c.SetTTL(time.Minute)
	assert.Equal(t, &cache{ttl: time.Minute}, c)
}

func TestSetCacheStorage(t *testing.T) {
	storage := NewInMemory()
	c := new(cache)

	SetCacheStorage(storage)(c)
	assert.Equal(t, &cache{storage: storage}, c)
}

func TestSetCacheTTL(t *testing.T) {
	c := new(cache)

	SetCacheTTL(time.Minute)(c)
	assert.Equal(t, &cache{ttl: time.Minute}, c)
}

func Test_newCache(t *testing.T) {
	c := newCache(SetCacheTTL(time.Minute))
	assert.Equal(t, &cache{storage: NewInMemory(), ttl: time.Minute}, c)
}
