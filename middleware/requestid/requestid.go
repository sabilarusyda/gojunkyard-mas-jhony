package requestid

import (
	"context"
	"net/http"

	"devcode.xeemore.com/systech/gojunkyard/util"
)

type ctxK struct{}

// New ...
func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const header = "x-request-id"
			id := r.Header.Get(header)
			if len(id) == 0 {
				id = util.GenerateUUID()
			}
			w.Header().Set(header, id)
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxK{}, id)))
		})
	}
}

// GetFromRequest ...
func GetFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	return GetFromContext(r.Context())
}

// GetFromContext ...
func GetFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(ctxK{}).(string)
	return id
}
