package panic

import (
	"net/http"
	"runtime/debug"

	"devcode.xeemore.com/systech/gojunkyard/errors"
	"devcode.xeemore.com/systech/gojunkyard/http/httpresponse"

	"github.com/julienschmidt/httprouter"
)

type Reporter interface {
	ReportHTTPPanic(err interface{}, stacktrace []byte, r *http.Request) error
}

// InitHTTPRouterRecover is used to initiate panic recover middleware
func InitHTTPRouterRecover(rp Reporter) func(httprouter.Handle) httprouter.Handle {
	return func(h httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			defer func() {
				if err := recover(); err != nil {
					httpresponse.WithError(w, 500, errors.New(500, "LG001", "Internal Server Error", "Internal Server Error"))
					if rp != nil {
						rp.ReportHTTPPanic(err, debug.Stack(), r)
					}
				}
			}()
			h(w, r, ps)
		}
	}
}

// InitStdRecover is used to initiate panic recover middleware
func InitStdRecover(rp Reporter) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					httpresponse.WithError(w, 500, errors.New(500, "LG001", "Internal Server Error", "Internal Server Error"))
					if rp != nil {
						rp.ReportHTTPPanic(err, debug.Stack(), r)
					}
				}
			}()
			h.ServeHTTP(w, r)
		}
	}
}
