package newrelic

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

// NewHandler ...
func NewHandler(app newrelic.Application) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txn := app.StartTransaction(r.Method+" "+r.URL.Path, w, r)
			defer txn.End()

			r = newrelic.RequestWithTransactionContext(r, txn)

			h.ServeHTTP(txn, r)
		})
	}
}

// NewHandlerFunc ...
func NewHandlerFunc(app newrelic.Application) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return NewHandler(app)(h).(http.HandlerFunc)
	}
}
