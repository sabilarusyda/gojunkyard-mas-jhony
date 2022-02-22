package slack

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Slack struct {
	app         string
	hookURL     string
	httpClient  iHTTPClient
	payloadPool *sync.Pool
	buffPool    *sync.Pool
}

type slackPayload struct {
	Text        string             `json:"text"`
	Attachments [1]slackAttachment `json:"attachments"`
}

type slackAttachment struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Color string `json:"color"`
}

type iHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ErrStatusNotOK is used while http response status is not 200
var ErrStatusNotOK = errors.New("Status not ok")

// mock time.Now
var now = time.Now

// NewSlackReporter is used to initiate slack reporter
func NewSlackReporter(appName, hookURL string) *Slack {
	return &Slack{
		app:         appName,
		hookURL:     hookURL,
		httpClient:  &http.Client{Timeout: time.Second * 5},
		payloadPool: &sync.Pool{New: func() interface{} { return new(slackPayload) }},
		buffPool:    &sync.Pool{New: func() interface{} { return bytes.NewBuffer(make([]byte, 0, 1500)) }},
	}
}

func (s *Slack) Debug(v ...interface{})                   {}
func (s *Slack) Debugf(format string, v ...interface{})   {}
func (s *Slack) Debugln(v ...interface{})                 {}
func (s *Slack) Info(v ...interface{})                    {}
func (s *Slack) Infof(format string, v ...interface{})    {}
func (s *Slack) Infoln(v ...interface{})                  {}
func (s *Slack) Warning(v ...interface{})                 {}
func (s *Slack) Warningf(format string, v ...interface{}) {}
func (s *Slack) Warningln(v ...interface{})               {}
func (s *Slack) Error(v ...interface{})                   {}
func (s *Slack) Errorf(format string, v ...interface{})   {}
func (s *Slack) Errorln(v ...interface{})                 {}

// ReportPanic is used to send the panic message to slack
func (s *Slack) ReportPanic(err interface{}, stacktrace []byte) error {
	payload := s.payloadPool.Get().(*slackPayload)
	payload.Text = fmt.Sprintf("<@here> *[PANIC!!!] APP:* `%s` *| TIME:* `%s`", s.app, now())
	payload.Attachments[0].Title = "Stacktrace:"
	payload.Attachments[0].Color = "danger"
	payload.Attachments[0].Text = s.format(err, string(stacktrace))
	defer s.payloadPool.Put(payload)
	return s.send(payload)
}

// ReportHTTPPanic is used to send the panic message to slack
func (s *Slack) ReportHTTPPanic(err interface{}, stacktrace []byte, _ *http.Request) error {
	return s.ReportPanic(err, stacktrace)
}

func (s *Slack) format(vs ...interface{}) string {
	// step 1. get buffer from pool
	buf := s.buffPool.Get().(*bytes.Buffer)
	buf.Reset()

	// step 2. write the payload
	buf.WriteString("```\n")
	for _, v := range vs {
		fmt.Fprintln(buf, v)
	}
	buf.WriteString("```")

	// step 3. put the buffer t
	defer s.buffPool.Put(buf)
	return buf.String()
}

func (s *Slack) send(payload *slackPayload) error {
	// step 1. encode the payload to json format
	byt, _ := jsoniter.ConfigFastest.Marshal(payload)

	// step 2. create a request
	req, _ := http.NewRequest(http.MethodPost, s.hookURL, bytes.NewBuffer(byt))

	// step 3. do the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// step 4. validate statusCode
	if resp.StatusCode != http.StatusOK {
		return ErrStatusNotOK
	}

	// step 5. return err nil
	return nil
}
