package valkyrie

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/json-iterator/go"
)

type core struct {
	client iHTTPClient
	h2p    map[string]int64
	i2p    map[string]int64
	option *Option
	mux    sync.RWMutex
	ts     time.Time
}

type iHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Option is used to initialize this middleware
// URL: endpoint of hostmap
type Option struct {
	URL           string
	CacheDuration time.Duration
}

type (
	hostmapResponse struct {
		Data []struct {
			ProjectID  int64 `json:"id"`
			Attributes struct {
				Hosts      []string `json:"host"`
				Identifier string   `json:"identifier"`
			}
		} `json:"data"`
	}
)

type ctxValue string

// PID means projectID. This is used to get project id from context
const (
	PID ctxValue = "valkyrie-pid"
)

var logFatalf = log.Fatalf

func construct(option *Option) *core {
	examineOption(option)
	c := &core{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		option: option,
	}

	log.Print("[valkyrie] initializing hostmap")
	for c.loadHostmap() != nil {
		log.Print("[valkyrie] failed to load hostmap")
		time.Sleep(3 * time.Second)
	}
	log.Print("[valkyrie] hostmap successfully loaded")

	return c
}

func examineOption(option *Option) {
	if option == nil {
		logFatalf("Failed to initialize valkyrie. option cannot be nil")
		return
	}
	if len(option.URL) == 0 {
		logFatalf("Failed to initialize valkyrie. option.url cannot be empty")
		return
	}
	if _, err := url.ParseRequestURI(option.URL); err != nil {
		logFatalf("Failed to initialize valkyrie. option.url got invalid format. err: %s", err)
		return
	}
	if option.CacheDuration == 0 {
		option.CacheDuration = time.Minute
	}
}

func (c *core) loadHostmap() error {
	req, _ := http.NewRequest(http.MethodGet, c.option.URL, nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("[loadHostmap] returns code: " + strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}

	var hostmap hostmapResponse
	err = jsoniter.Unmarshal(body, &hostmap)
	if err != nil {
		return err
	}

	c.mux.Lock()
	c.i2p = make(map[string]int64, len(hostmap.Data))
	c.h2p = make(map[string]int64, len(hostmap.Data)*2)
	for _, v := range hostmap.Data {
		c.i2p[v.Attributes.Identifier] = v.ProjectID
		for _, v2 := range v.Attributes.Hosts {
			c.h2p[v2] = v.ProjectID
		}
	}
	c.ts = time.Now()
	c.mux.Unlock()

	return nil
}

func (c *core) authProject(r *http.Request) (pid int64, ok bool) {
	identifier := r.FormValue("app_id")
	if len(identifier) == 0 {
		return
	}
	pid, ok = c.geti2p(identifier)
	if ok {
		return
	}
	if c.renewHostmap() != nil {
		return
	}
	pid, ok = c.geti2p(identifier)
	return
}

func (c *core) authHostname(r *http.Request) (pid int64, ok bool) {
	hostname := r.Host
	if len(r.Header.Get("x-url")) > 0 {
		hostname = r.Header.Get("x-url")
	}
	pid, ok = c.geth2p(hostname)
	if ok {
		return
	}
	if c.renewHostmap() != nil {
		return
	}
	pid, ok = c.geth2p(hostname)
	return
}

func (c *core) geth2p(hostname string) (pid int64, ok bool) {
	c.mux.RLock()
	pid, ok = c.h2p[hostname]
	c.mux.RUnlock()
	return
}

func (c *core) geti2p(identifier string) (pid int64, ok bool) {
	c.mux.RLock()
	pid, ok = c.i2p[identifier]
	c.mux.RUnlock()
	return
}

func (c *core) renewHostmap() error {
	if expTime := c.ts.Add(c.option.CacheDuration); time.Now().Before(expTime) {
		return errors.New("cache is still valid")
	}
	return c.loadHostmap()
}
