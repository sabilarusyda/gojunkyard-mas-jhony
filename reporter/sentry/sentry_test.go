package sentry

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	_sentry "github.com/getsentry/raven-go"
	"github.com/stretchr/testify/assert"
)

func TestNewSentryReporter(t *testing.T) {
	// TEST 1
	option := &Option{
		AppName:     "devcode.xeemore.com/systech/gojunkyard",
		SampleRate:  1.0,
		Environment: "production",
		AppVersion:  "v1.0.8",
	}
	got := NewSentryReporter(option)

	client, _ := _sentry.NewWithTags(option.DSN, map[string]string{
		"app_name": option.AppName,
	})
	client.SetSampleRate(option.SampleRate)
	client.SetEnvironment(option.Environment)
	client.SetRelease(option.AppVersion)

	assert.Equal(t, option, got.option)
	assert.Equal(t, client.Tags, got.client.Tags)

	// TEST 2
	option = &Option{
		AppName: "devcode.xeemore.com/systech/gojunkyard",
	}
	got = NewSentryReporter(option)

	client, _ = _sentry.NewWithTags(option.DSN, map[string]string{
		"app_name": option.AppName,
	})

	assert.Equal(t, option, got.option)
	assert.Equal(t, client.Tags, got.client.Tags)
}

func TestSentry_DEBUG_INFO_WARNING_ERROR(t *testing.T) {
	s := NewSentryReporter(&Option{AppName: "devcode.xeemore.com/systech/gojunkyard", Level: DEBUG})
	s.Debug()
	s.Debugf("")
	s.Debugln()
	s.Info()
	s.Infof("")
	s.Infoln()
	s.Warning()
	s.Warningf("")
	s.Warningln()
	s.Error()
	s.Errorf("")
	s.Errorln()
}

func TestSentry_generatePanicPacket(t *testing.T) {
	s := NewSentryReporter(&Option{AppName: "devcode.xeemore.com/systech/gojunkyard", Level: DEBUG})

	// case 1, err nil
	assert.Nil(t, s.generatePanicPacket(nil))

	// case 2, error
	want := _sentry.NewPacket("===ERROR===", _sentry.NewException(errors.New("===ERROR==="), _sentry.NewStacktrace(3, 3, nil)))
	want.Level = _sentry.FATAL
	got := s.generatePanicPacket(errors.New("===ERROR==="))
	assert.Equal(t, want, got)

	// case 3, non error
	want = _sentry.NewPacket("123", _sentry.NewException(errors.New("123"), _sentry.NewStacktrace(3, 3, nil)))
	want.Level = _sentry.FATAL
	got = s.generatePanicPacket(123)
	assert.Equal(t, want, got)
}

func TestSentry_ReportPanic(t *testing.T) {
	s := NewSentryReporter(&Option{AppName: "devcode.xeemore.com/systech/gojunkyard", Level: DEBUG})
	err := s.ReportPanic(errors.New("===ERROR==="), []byte("===STACKTRACE==="))
	assert.Nil(t, err)
}

func TestSentry_ReportHTTPPanic(t *testing.T) {
	s := NewSentryReporter(&Option{AppName: "devcode.xeemore.com/systech/gojunkyard", Level: DEBUG})
	// case 1. use report panic function
	err := s.ReportHTTPPanic(errors.New("===ERROR==="), []byte("===STACKTRACE==="), nil)
	assert.Nil(t, err)
	// case 2. use http panic function
	err = s.ReportHTTPPanic(errors.New("===ERROR==="), []byte("===STACKTRACE==="), httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Nil(t, err)
}
