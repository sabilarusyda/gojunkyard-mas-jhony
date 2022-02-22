package http

import (
	"fmt"
	"net/http"
	"time"
)

type HTTP struct {
	path   string
	name   string
	client interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

// New ...
func New(path string) *HTTP {
	return &HTTP{path: path}
}

// SetName ...
func (h *HTTP) SetName(name string) {
	h.name = name
}

// Name ...
func (h *HTTP) Name() string {
	if len(h.name) == 0 {
		return "HTTP"
	}
	return h.name
}

// Check ...
func (h *HTTP) Check() error {
	if h.client == nil {
		h.client = newClient()
	}

	req, err := http.NewRequest(http.MethodGet, h.path, nil)
	if err != nil {
		return err
	}

	res, err := h.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("Status Code %d", res.StatusCode)
	}

	return nil
}

func newClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
		},
	}
}
