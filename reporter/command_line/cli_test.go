package cli

import (
	"net/http"
	"os"
	"testing"

	"devcode.xeemore.com/systech/gojunkyard/reporter/writer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type _reporter struct {
	mock.Mock
}

func (r *_reporter) SetCallDepth(calldepth int) {
	r.Called(calldepth)
}

func (r *_reporter) SetFlags(flag int) {
	r.Called(flag)
}

func (r *_reporter) Debug(v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Debugf(format string, v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Debugln(v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Info(v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Infof(format string, v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Infoln(v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Warning(v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Warningf(format string, v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Warningln(v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Error(v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Errorf(format string, v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) Errorln(v ...interface{}) {
	r.Called(nil)
}

func (r *_reporter) ReportPanic(err interface{}, stacktrace []byte) error {
	return r.Called(nil, nil).Error(0)
}

func (r *_reporter) ReportHTTPPanic(err interface{}, stacktrace []byte, req *http.Request) error {
	return r.Called(nil, nil, nil).Error(0)
}

func TestNewCliReporter(t *testing.T) {
	const appName = "devcode.xeemore.com/systech/gojunkyardjunkyard"
	const level = 0
	const calldepth = 4

	stdout := writer.NewWriterReporter(appName, level, os.Stdout)
	stderr := writer.NewWriterReporter(appName, level, os.Stderr)

	var (
		want = &CliReporter{
			stdout: stdout,
			stderr: stderr,
		}
		got = NewCliReporter(appName, level)
	)
	assert.Equal(t, want, got)
}

func TestCliReporter_SetCallDepth(t *testing.T) {
	stdout := new(_reporter)
	stderr := new(_reporter)
	reporter := NewCliReporter("", 1)
	reporter.stdout = stdout
	reporter.stderr = stderr

	stdout.On("SetCallDepth", 3).Once()
	stderr.On("SetCallDepth", 3).Once()
	reporter.SetCallDepth(3)
}

func TestCliReporter_SetFlags(t *testing.T) {
	stdout := new(_reporter)
	stderr := new(_reporter)
	reporter := NewCliReporter("", 1)
	reporter.stdout = stdout
	reporter.stderr = stderr

	stdout.On("SetFlags", LUTC|Lshortfile).Once()
	stderr.On("SetFlags", LUTC|Lshortfile).Once()
	reporter.SetFlags(LUTC | Lshortfile)
}

func TestCliReporter_Debug(t *testing.T) {
	const (
		appName = "devcode.xeemore.com/systech/gojunkyardjunkyard"
		level   = 1
	)
	var (
		stdout   = new(_reporter)
		stderr   = new(_reporter)
		reporter = NewCliReporter(appName, level)
	)
	reporter.stderr = stderr
	reporter.stdout = stdout

	stdout.On("Debug", nil)
	stdout.On("Debugf", nil)
	stdout.On("Debugln", nil)
	reporter.Debug()
	reporter.Debugf("")
	reporter.Debugln()

	stdout.On("Info", nil)
	stdout.On("Infof", nil)
	stdout.On("Infoln", nil)
	reporter.Info()
	reporter.Infof("")
	reporter.Infoln()

	stdout.On("Warning", nil)
	stdout.On("Warningf", nil)
	stdout.On("Warningln", nil)
	reporter.Warning()
	reporter.Warningf("")
	reporter.Warningln()

	stderr.On("Error", nil)
	stderr.On("Errorf", nil)
	stderr.On("Errorln", nil)
	reporter.Error()
	reporter.Errorf("")
	reporter.Errorln()

	stderr.On("ReportPanic", nil, nil).Return(nil)
	stderr.On("ReportHTTPPanic", nil, nil, nil).Return(nil)
	reporter.ReportPanic(nil, nil)
	reporter.ReportHTTPPanic(nil, nil, nil)
}
