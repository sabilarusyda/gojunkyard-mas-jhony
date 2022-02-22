package panic

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type _reporter struct {
	mock.Mock
}

func (r *_reporter) ReportHTTPPanic(err interface{}, stacktrace []byte, req *http.Request) error {
	return r.Called(err, nil, req).Error(0)
}

func TestHTTPRouterRecoverWithReporter(t *testing.T) {
	// step 1. prepare wanted result
	const (
		wantBody   = `{"errors":[{"status":"500","code":"LG001","title":"Internal Server Error","detail":"Internal Server Error"}]}`
		wantStatus = 500
		panicmsg   = "===TEST==="
	)

	// step 2. prepare testing variable
	var (
		reporter = new(_reporter)
		mw       = InitHTTPRouterRecover(reporter)
		f        = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			panic(panicmsg)
		}
		r = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		w = httptest.NewRecorder()
	)

	// step 3. prepare mocking object (we dont need to test stacktrace because is will change for every execution)
	reporter.On("ReportHTTPPanic", panicmsg, nil, r).Return(nil)

	// step 4. call the function wanted to test
	mw(f)(w, r, nil)

	// step 5. compare the function result to wanted result
	assert.Equal(t, wantBody, strings.TrimSpace(w.Body.String()))
	assert.Equal(t, wantStatus, w.Code)
}

func TestHTTPRouterRecoverWithoutReporter(t *testing.T) {
	// step 1. prepare wanted result
	const (
		wantBody   = `{"errors":[{"status":"500","code":"LG001","title":"Internal Server Error","detail":"Internal Server Error"}]}`
		wantStatus = 500
		panicmsg   = "===TEST==="
	)

	// step 2. prepare testing variable
	var (
		mw = InitHTTPRouterRecover(nil)
		f  = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			panic(panicmsg)
		}
		r = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		w = httptest.NewRecorder()
	)

	// step 3. call the function wanted to test
	mw(f)(w, r, nil)

	// step 4. compare the function result to wanted result
	assert.Equal(t, wantBody, strings.TrimSpace(w.Body.String()))
	assert.Equal(t, wantStatus, w.Code)
}

func TestStdRecoverWithReporter(t *testing.T) {
	// step 1. prepare wanted result
	const (
		wantBody   = `{"errors":[{"status":"500","code":"LG001","title":"Internal Server Error","detail":"Internal Server Error"}]}`
		wantStatus = 500
		panicmsg   = "===TEST==="
	)

	// step 2. prepare testing variable
	var (
		reporter = new(_reporter)
		mw       = InitStdRecover(reporter)
		f        = func(w http.ResponseWriter, r *http.Request) {
			panic(panicmsg)
		}
		r = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		w = httptest.NewRecorder()
	)

	// step 3. prepare mocking object (we dont need to test stacktrace because is will change for every execution)
	reporter.On("ReportHTTPPanic", panicmsg, nil, r).Return(nil)

	// step 4. call the function wanted to test
	mw(f)(w, r)

	// step 5. compare the function result to wanted result
	assert.Equal(t, wantBody, strings.TrimSpace(w.Body.String()))
	assert.Equal(t, wantStatus, w.Code)
}

func TestStdRecoverWithoutReporter(t *testing.T) {
	// step 1. prepare wanted result
	const (
		wantBody   = `{"errors":[{"status":"500","code":"LG001","title":"Internal Server Error","detail":"Internal Server Error"}]}`
		wantStatus = 500
		panicmsg   = "===TEST==="
	)

	// step 2. prepare testing variable
	var (
		mw = InitStdRecover(nil)
		f  = func(w http.ResponseWriter, r *http.Request) {
			panic(panicmsg)
		}
		r = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		w = httptest.NewRecorder()
	)

	// step 3. call the function wanted to test
	mw(f)(w, r)

	// step 4. compare the function result to wanted result
	assert.Equal(t, wantBody, strings.TrimSpace(w.Body.String()))
	assert.Equal(t, wantStatus, w.Code)
}
