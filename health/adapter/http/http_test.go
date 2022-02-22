package http

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type _httpClient struct {
	mock.Mock
}

func (hc *_httpClient) Do(req *http.Request) (*http.Response, error) {
	args := hc.Called()
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNew(t *testing.T) {
	const method = http.MethodGet
	const path = "http://127.0.0.1:3000/ping"

	checker := New(path)
	assert.Equal(t, &HTTP{path: path}, checker)
}

func TestHTTP_GetSetName(t *testing.T) {
	const method = http.MethodGet
	const path = "http://127.0.0.1:3000/ping"

	http := New(path)
	assert.Equal(t, "HTTP", http.Name())

	const name = "HTTP " + path
	http.SetName(name)
	assert.Equal(t, name, http.Name())
}

func Test_newClient(t *testing.T) {
	assert.Equal(t, &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
		},
	}, newClient())
}

func TestHTTP_Check_NilHTTPClient(t *testing.T) {
	const method = http.MethodGet
	const path = "O(*AS)D(A*SD)(AS*D)(*ASD)(*AS)D("

	var response *http.Response
	hc := new(_httpClient)
	hc.On("Do").Once().Return(response, http.ErrServerClosed)

	checker := New(path)
	assert.Error(t, checker.Check())
}

func TestHTTP_Check_ErrConn(t *testing.T) {
	const method = http.MethodGet
	const path = "http://127.0.0.1:3000/ping"

	var response *http.Response
	hc := new(_httpClient)
	hc.On("Do").Once().Return(response, http.ErrServerClosed)

	checker := New(path)
	checker.client = hc
	assert.Equal(t, http.ErrServerClosed, checker.Check())
}

func TestHTTP_Check_ServiceUnavailable(t *testing.T) {
	const path = "http://127.0.0.1:3000/ping"
	const code = http.StatusServiceUnavailable

	hc := new(_httpClient)
	hc.On("Do").Once().Return(&http.Response{
		Status:     http.StatusText(code),
		StatusCode: code,
	}, nil)

	checker := New(path)
	checker.client = hc
	assert.Equal(t, fmt.Errorf("Status Code %d", code), checker.Check())
}

func TestHTTP_Check_ServiceSuccess(t *testing.T) {
	const path = "http://127.0.0.1:3000/ping"
	const code = http.StatusOK

	hc := new(_httpClient)
	hc.On("Do").Once().Return(&http.Response{
		Status:     http.StatusText(code),
		StatusCode: code,
	}, nil)

	checker := New(path)
	checker.client = hc
	assert.Nil(t, checker.Check())
}

func TestHTTP_Name(t *testing.T) {
	type fields struct {
		path   string
		name   string
		client interface {
			Do(req *http.Request) (*http.Response, error)
		}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTP{
				path:   tt.fields.path,
				name:   tt.fields.name,
				client: tt.fields.client,
			}
			if got := h.Name(); got != tt.want {
				t.Errorf("HTTP.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
