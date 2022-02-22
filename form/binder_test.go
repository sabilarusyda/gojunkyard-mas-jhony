package form

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type user struct {
	ID        int64     `form:"id"`
	Name      string    `form:"name" mod:"trim"`
	Email     string    `form:"email" mod:"trim"`
	Gender    string    `form:"gender" validate:"required,eq=m|eq=f"`
	BirthDate time.Time `form:"bdate"`
}

func TestBindFlagNilRequest(t *testing.T) {
	const payload = "@!^#%&%!@#!@&#"
	var (
		got user
		r   *http.Request
	)

	err := BindFlag(Bnone, &got, r)
	assert.NotNil(t, err, "error must be exist due to nil request")
}

func TestBindFlagFailedParseForm(t *testing.T) {
	const payload = "@!^#%&%!@#!@&#"
	var (
		got user
		r   = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := BindFlag(Bnone, &got, r)
	assert.NotNil(t, err, "error must be exist")
}

func BenchmarkBindFlagFailedParseForm(b *testing.B) {
	const payload = "@!^#%&%!@#!@&#"
	var (
		r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var got user
		BindFlag(Bnone, &got, r)
	}
}

func TestBindFlagFailedDecode(t *testing.T) {
	var (
		got     user
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := BindFlag(Bnone, got, r)
	assert.NotNil(t, err, "error must be exist")
}

func BenchmarkBindFlagFailedDecode(b *testing.B) {
	var (
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var got user
		BindFlag(Bnone, &got, r)
	}
}

func TestBindFlagFailedDecodeJSON(t *testing.T) {
	var (
		got     user
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := BindFlag(Bnone, got, r)
	assert.NotNil(t, err, "error must be exist")
}

func TestBindFlagFailedUnsupportedContentType(t *testing.T) {
	var (
		got     user
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := BindFlag(Bnone, got, r)
	assert.Equal(t, ErrUnsupportedContentType, err)
}

func TestBindFlagSuccessMethodGet(t *testing.T) {
	var (
		got  user
		want = user{ID: 1408, Name: "Risal Falah"}
		r    = httptest.NewRequest(http.MethodGet, "/", nil)
	)
	q := r.URL.Query()
	q.Add("id", "1408")
	q.Add("name", "Risal Falah")
	r.URL.RawQuery = q.Encode()

	err := BindFlag(Bnone, &got, r)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, want.BirthDate.Format(time.RFC3339), got.BirthDate.Format(time.RFC3339))
	got.BirthDate = want.BirthDate // bypass from checking
	assert.Equal(t, want, got)
}

func BenchmarkBindFlagFailedDecodeJSON(b *testing.B) {
	var (
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var got user
		BindFlag(Bnone, &got, r)
	}
}

func TestBindFlagSuccessFilter(t *testing.T) {
	var (
		got     user
		want    = user{ID: 1408, Name: "Risal Falah"}
		payload = url.Values{"id": {"1408"}, "name": {" Risal Falah "}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := BindFlag(Bfilter, &got, r)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, want.BirthDate.Format(time.RFC3339), got.BirthDate.Format(time.RFC3339))
	got.BirthDate = want.BirthDate // bypass from checking
	assert.Equal(t, want, got)
}

func TestBindFlagFailedValidation(t *testing.T) {
	var (
		got     user
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}, "gender": {"x"}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := BindFlag(Bfilter|Bvalidate, &got, r)
	assert.NotNil(t, err, "error must be exist")
}

func TestBindFlagSuccessValidation(t *testing.T) {
	var (
		got     user
		dt      = time.Date(2018, time.December, 1, 1, 1, 1, 0, time.Local)
		want    = user{ID: 1408, Name: "Risal Falah", Gender: "m", BirthDate: dt}
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}, "gender": {"m"}, "bdate": {dt.Format(time.RFC3339)}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := BindFlag(Bfilter|Bvalidate, &got, r)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, want.BirthDate.Format(time.RFC3339), got.BirthDate.Format(time.RFC3339))
	got.BirthDate = want.BirthDate // bypass from checking
	assert.Equal(t, want, got)
}

func BenchmarkBindFlagSuccessValidation(b *testing.B) {
	var (
		dt      = time.Date(2018, time.December, 1, 1, 1, 1, 0, time.Local)
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}, "gender": {"m"}, "bdate": {dt.Format(time.RFC3339)}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var got user
		err := BindFlag(Bnone, &got, r)
		if err != nil {
			b.Error("error must be nil")
		}
	}
}

func TestBindFlagSuccessParseForm(t *testing.T) {
	var (
		got     user
		want    = user{ID: 1408, Name: "Risal Falah"}
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := BindFlag(Bnone, &got, r)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, want, got)
}

func BenchmarkBindFlagSuccessParseForm(b *testing.B) {
	var (
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var got user
		err := BindFlag(Bnone, &got, r)
		if err != nil {
			b.Error("error must be nil")
		}
	}
}

func TestBindSuccessValidation(t *testing.T) {
	var (
		got     user
		dt      = time.Date(2018, time.December, 1, 1, 1, 1, 0, time.Local)
		want    = user{ID: 1408, Name: "Risal Falah", Gender: "m", BirthDate: dt}
		payload = url.Values{"id": {"1408"}, "name": {"Risal Falah"}, "gender": {"m"}, "bdate": {dt.Format(time.RFC3339)}}.Encode()
		r       = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := Bind(&got, r)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, want.BirthDate.Format(time.RFC3339), got.BirthDate.Format(time.RFC3339))
	got.BirthDate = want.BirthDate // bypass from checking
	assert.Equal(t, want, got)
}

func TestBindJSONSuccessValidation(t *testing.T) {
	var (
		got     user
		want    = user{ID: 1408, Name: "Risal Falah", Gender: "m", BirthDate: time.Date(2018, time.January, 1, 1, 2, 2, 2, time.UTC)}
		payload = []byte(`
		{
			"id": 1408,
			"name": "Risal Falah",
			"gender": "m",
			"bdate": "2018-01-01T01:02:02.000000002Z"
		}
		`)
		r = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
	)
	r.Header.Add("Content-Type", "application/json;charset=UTF-8")
	r.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	err := Bind(&got, r)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, want, got)
}
