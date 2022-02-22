package health

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"devcode.xeemore.com/systech/gojunkyard/webserver"

	"devcode.xeemore.com/systech/gojunkyard/reporter"
	"devcode.xeemore.com/systech/gojunkyard/reporter/nop"
)

// Checker is interface of health check adapter
type Checker interface {
	Name() string
	Check() error
}

// Health is object that handle liveness (ping) and readiness (healthz)
type Health struct {
	mux      sync.RWMutex
	checker  []Checker
	ready    bool
	server   webserver.Server
	reporter reporter.Reporter
}

// New returns Health object
func New() *Health {
	health := &Health{
		checker:  make([]Checker, 0),
		reporter: nop.NewNopReporter(),
		server: webserver.New(&webserver.Options{
			ListenAddress:   ":54322",
			MaxConnections:  5,
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			GracefulTimeout: 5 * time.Second,
		}),
	}
	health.initrouter()
	return health
}

func (h *Health) initrouter() {
	router := h.server.Router()
	router.GET("/ping", h.ping)
	router.GET("/healthz", h.healthz)
}

// Register ...
func (h *Health) Register(c ...Checker) {
	h.mux.Lock()
	h.checker = append(h.checker, c...)
	h.mux.Unlock()
}

// Run ...
func (h *Health) Run() chan error {
	return h.server.Run()
}

func (h *Health) ping(w http.ResponseWriter, _ *http.Request) {
	const pong = "PONG"
	w.Write([]byte(pong))
	h.reporter.Infof("PING: [%s]\n", pong)
}

// SetReadiness ...
func (h *Health) SetReadiness(ready bool) {
	h.mux.Lock()
	h.ready = ready
	h.mux.Unlock()
}

// GetReadiness ...
func (h *Health) GetReadiness() bool {
	h.mux.RLock()
	defer h.mux.RUnlock()
	return h.ready
}

// SetReporter ...
func (h *Health) SetReporter(reporter reporter.Reporter) {
	h.reporter = reporter
}

func (h *Health) healthz(w http.ResponseWriter, _ *http.Request) {
	// condition 1. if app has not set to be ready then return service unavailable
	if !h.GetReadiness() {
		const msg = "Not ready to check. App is trying to up"
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(msg))
		h.reporter.Warningf("Health: %s\n", msg)
		return
	}

	// condition 2. app health has been ready to be checked
	var (
		code = http.StatusOK
		buff = bytes.NewBuffer(make([]byte, 0, len(h.checker)*20))
	)

	for _, v := range h.checker {
		err := v.Check()
		if err == nil {
			fmt.Fprintf(buff, "%s: [OK]\n", v.Name())
			continue
		}

		fmt.Fprintf(buff, "%s: [err: %s]\n", v.Name(), err)
		code = http.StatusServiceUnavailable
	}

	// get the data
	var byt = buff.Bytes()

	var b string
	bh := (*reflect.StringHeader)(unsafe.Pointer(&b))
	bh.Data = (*reflect.StringHeader)(unsafe.Pointer(&byt)).Data
	bh.Len = len(byt)

	/* write the response */
	w.WriteHeader(code)
	w.Write(byt)

	/* log the health */
	reportf := h.reporter.Infof
	if code != http.StatusOK {
		reportf = h.reporter.Errorf
	}
	reportf("Health: \n%s\n", b)
}

// RunGraceful run the webserver with blocking
func (h *Health) RunGraceful() error {
	h.SetReadiness(true)
	defer h.SetReadiness(false)

	return h.server.RunGraceful()
}

// Stop terminate the server gracefully
func (h *Health) Stop() error {
	if h.server == nil {
		return nil
	}

	defer h.SetReadiness(false)
	return h.server.Stop()
}
