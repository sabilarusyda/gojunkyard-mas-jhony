package slack

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type _httpError struct{}

func (s *_httpError) Do(req *http.Request) (*http.Response, error) {
	return nil, http.ErrServerClosed
}

type _httpInternalServerError struct{}

func (s *_httpInternalServerError) Do(req *http.Request) (*http.Response, error) {
	const code = http.StatusInternalServerError
	var text = http.StatusText(code)
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(text)),
		StatusCode: code,
		Status:     text,
	}, nil
}

type _httpSuccess struct{}

func (s *_httpSuccess) Do(req *http.Request) (*http.Response, error) {
	const code = http.StatusOK
	var text = http.StatusText(code)
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(text)),
		StatusCode: code,
		Status:     text,
	}, nil
}

func TestNewSlackReporter(t *testing.T) {
	const appName = "supersoccer"
	const hookURL = "https://hooks.slack.com/services/ABCDEFG"
	var got = NewSlackReporter(appName, hookURL)
	assert.Equal(t, appName, got.app)
	assert.Equal(t, hookURL, got.hookURL)
	assert.Equal(t, &http.Client{Timeout: time.Second * 5}, got.httpClient)
	assert.Equal(t, new(slackPayload), got.payloadPool.Get().(*slackPayload))
	assert.Equal(t, bytes.NewBuffer(make([]byte, 0, 1500)), got.buffPool.Get().(*bytes.Buffer))
}

func TestSlack_format(t *testing.T) {
	const want = "```\n1\n2\n```"
	got := NewSlackReporter("", "").format(1, 2)
	assert.Equal(t, want, got)
}

func TestSlack_ReportHTTPPanic(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	reporter := NewSlackReporter("", "")
	reporter.httpClient = new(_httpSuccess)
	err := reporter.ReportHTTPPanic("TEST", []byte("===STACKTRACE==="), request)
	assert.Nil(t, err)

	reporter.httpClient = new(_httpInternalServerError)
	err = reporter.ReportHTTPPanic("TEST", []byte("===STACKTRACE==="), request)
	assert.NotNil(t, err)

	reporter.httpClient = new(_httpError)
	err = reporter.ReportHTTPPanic("TEST", []byte("===STACKTRACE==="), request)
	assert.NotNil(t, err)
}
