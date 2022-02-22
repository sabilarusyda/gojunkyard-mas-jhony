package health

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"devcode.xeemore.com/systech/gojunkyard/reporter"
	"devcode.xeemore.com/systech/gojunkyard/webserver"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type _reporter struct {
	mock.Mock
}

func (r *_reporter) Debug(v ...interface{}) {
	r.Called(v...)
}

func (r *_reporter) Debugf(format string, v ...interface{}) {
	iface := make([]interface{}, 0, len(v)+1)
	iface = append(iface, format)
	iface = append(iface, v...)
	r.Called(iface...)
}

func (r *_reporter) Debugln(v ...interface{}) {
	r.Called(v...)
}

func (r *_reporter) Info(v ...interface{}) {
	r.Called(v...)
}

func (r *_reporter) Infof(format string, v ...interface{}) {
	iface := make([]interface{}, 0, len(v)+1)
	iface = append(iface, format)
	iface = append(iface, v...)
	r.Called(iface...)
}

func (r *_reporter) Infoln(v ...interface{}) {
	r.Called(v...)
}

func (r *_reporter) Warning(v ...interface{}) {
	r.Called(v...)
}

func (r *_reporter) Warningf(format string, v ...interface{}) {
	iface := make([]interface{}, 0, len(v)+1)
	iface = append(iface, format)
	iface = append(iface, v...)
	r.Called(iface...)
}

func (r *_reporter) Warningln(v ...interface{}) {
	r.Called(v...)
}

func (r *_reporter) Error(v ...interface{}) {
	r.Called(v...)
}

func (r *_reporter) Errorf(format string, v ...interface{}) {
	iface := make([]interface{}, 0, len(v)+1)
	iface = append(iface, format)
	iface = append(iface, v...)
	r.Called(iface...)
}

func (r *_reporter) Errorln(v ...interface{}) {
	r.Called(v...)
}

func (r *_reporter) ReportPanic(err interface{}, stacktrace []byte) error {
	args := r.Called(err, stacktrace)
	return args.Error(0)
}

func (r *_reporter) ReportHTTPPanic(err interface{}, stacktrace []byte, req *http.Request) error {
	return r.Called(err, stacktrace, req).Error(0)
}

type _checker struct {
	mock.Mock
}

func (c *_checker) Name() string {
	return c.Called().Get(0).(string)
}

func (c *_checker) Check() error {
	return c.Called().Error(0)
}

func TestHealth_GetSetReadiness(t *testing.T) {
	health := New()

	// default: false
	assert.False(t, health.GetReadiness())

	// change to true
	health.SetReadiness(true)
	assert.True(t, health.GetReadiness())

	health.SetReadiness(false)
	assert.False(t, health.GetReadiness())
}

func TestHealth_ping(t *testing.T) {
	reporter := new(_reporter)
	reporter.On("Infof", "PING: [%s]\n", "PONG").Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/ping", nil)

	health := New()
	health.SetReporter(reporter)
	health.ping(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "PONG", w.Body.String())
}

func TestHealth_healthz(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name   string
		health *Health
		args   args
		want   want
	}{
		{
			name: "not ready",
			health: &Health{
				ready: false,
				reporter: func() reporter.Reporter {
					reporter := new(_reporter)
					reporter.On("Warningf", "Health: %s\n", "Not ready to check. App is trying to up")
					return reporter
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/healthz", nil),
			},
			want: want{
				body: "Not ready to check. App is trying to up",
				code: http.StatusServiceUnavailable,
			},
		},
		{
			name: "exist error",
			health: &Health{
				ready: true,
				checker: func() []Checker {
					checker := new(_checker)
					checker.On("Name").Return("HTTP_CLIENT 127.0.0.1:3000")
					checker.On("Check").Return(errors.New("Cannot connect to server 127.0.0.1:3000"))
					return []Checker{checker}
				}(),
				reporter: func() reporter.Reporter {
					reporter := new(_reporter)
					reporter.On("Errorf", "Health: \n%s\n", "HTTP_CLIENT 127.0.0.1:3000: [err: Cannot connect to server 127.0.0.1:3000]\n")
					return reporter
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/healthz", nil),
			},
			want: want{
				body: "HTTP_CLIENT 127.0.0.1:3000: [err: Cannot connect to server 127.0.0.1:3000]\n",
				code: http.StatusServiceUnavailable,
			},
		},
		{
			name: "success",
			health: &Health{
				ready: true,
				checker: func() []Checker {
					var err error
					checker := new(_checker)
					checker.On("Name").Return("HTTP_CLIENT 127.0.0.1:3000")
					checker.On("Check").Return(err)
					return []Checker{checker}
				}(),
				reporter: func() reporter.Reporter {
					reporter := new(_reporter)
					reporter.On("Infof", "Health: \n%s\n", "HTTP_CLIENT 127.0.0.1:3000: [OK]\n")
					return reporter
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/healthz", nil),
			},
			want: want{
				body: "HTTP_CLIENT 127.0.0.1:3000: [OK]\n",
				code: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.health.healthz(tt.args.w, tt.args.r)
			assert.Equal(t, tt.want.body, tt.args.w.Body.String())
			assert.Equal(t, tt.want.code, tt.args.w.Code)
		})
	}
}

func TestHealth_Stop(t *testing.T) {
	type want struct {
		err   error
		ready bool
	}
	tests := []struct {
		name   string
		health *Health
		want   want
	}{
		{
			name:   "nil server",
			health: new(Health),
			want:   want{err: nil},
		},
		{
			name: "success",
			health: &Health{
				ready:  true,
				server: webserver.New(&webserver.Options{}),
			},
			want: want{err: nil, ready: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.err, tt.health.Stop())
			assert.Equal(t, tt.want.ready, tt.health.ready)
		})
	}
}

func TestHealth_Register(t *testing.T) {
	var (
		health  = New()
		checker = new(_checker)
	)

	assert.Len(t, health.checker, 0)

	health.Register(checker)
	assert.Len(t, health.checker, 1)

	health.Register(checker, checker, checker, checker)
	assert.Len(t, health.checker, 5)

	for _, v := range health.checker {
		assert.Equal(t, checker, v)
	}
}
