package aggregator

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"devcode.xeemore.com/systech/gojunkyard/reporter"

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

func TestNewAggregator(t *testing.T) {
	var (
		rep  = new(_reporter)
		want = &Aggregator{rs: []reporter.Reporter{rep, rep, rep}}
		got  = NewAggregator(rep, rep, rep)
	)
	assert.Equal(t, want, got)
}

func TestDebug(t *testing.T) {
	const str = "===TEST==="
	const method = "Debug"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, str)
	aggregator.Debug(str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestDebugf(t *testing.T) {
	const str = "===TEST==="
	const method = "Debugf"
	const format = "hello %s"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, format, str)
	aggregator.Debugf(format, str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestDebugln(t *testing.T) {
	const str = "===TEST==="
	const method = "Debugln"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, str)
	aggregator.Debugln(str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestInfo(t *testing.T) {
	const str = "===TEST==="
	const method = "Info"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, str)
	aggregator.Info(str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestInfof(t *testing.T) {
	const str = "===TEST==="
	const method = "Infof"
	const format = "hello %s"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, format, str)
	aggregator.Infof(format, str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestInfoln(t *testing.T) {
	const str = "===TEST==="
	const method = "Infoln"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, str)
	aggregator.Infoln(str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestWarning(t *testing.T) {
	const str = "===TEST==="
	const method = "Warning"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, str)
	aggregator.Warning(str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestWarningln(t *testing.T) {
	const str = "===TEST==="
	const method = "Warningln"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, str)
	aggregator.Warningln(str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestWarningf(t *testing.T) {
	const str = "===TEST==="
	const method = "Warningf"
	const format = "hello %s"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, format, str)
	aggregator.Warningf(format, str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestError(t *testing.T) {
	const str = "===TEST==="
	const method = "Error"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, str)
	aggregator.Error(str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestErrorf(t *testing.T) {
	const str = "===TEST==="
	const method = "Errorf"
	const format = "hello %s"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, format, str)
	aggregator.Errorf(format, str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestErrorln(t *testing.T) {
	const str = "===TEST==="
	const method = "Errorln"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep)
	)

	rep.On(method, str)
	aggregator.Errorln(str)

	rep.AssertNumberOfCalls(t, method, 3)
	rep.AssertExpectations(t)
}

func TestAggregator_ReportPanic(t *testing.T) {
	const str = "===TEST==="
	const method = "ReportPanic"
	var (
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep, rep, rep)
		err        = errors.New(str)
		stacktrace = []byte(str)
	)

	rep.On(method, err, stacktrace).Return(nil)
	aggregator.ReportPanic(err, stacktrace)

	rep.AssertNumberOfCalls(t, method, 5)
	rep.AssertExpectations(t)
}

func TestAggregator_ReportHTTPPanic(t *testing.T) {
	const str = "===TEST==="
	const method = "ReportHTTPPanic"
	var (
		r          = httptest.NewRequest(http.MethodGet, "/", nil)
		rep        = new(_reporter)
		aggregator = NewAggregator(rep, rep, rep, rep, rep)
		err        = errors.New(str)
		stacktrace = []byte(str)
	)

	rep.On(method, err, stacktrace, r).Return(nil)
	aggregator.ReportHTTPPanic(err, stacktrace, r)

	rep.AssertNumberOfCalls(t, method, 5)
	rep.AssertExpectations(t)
}
