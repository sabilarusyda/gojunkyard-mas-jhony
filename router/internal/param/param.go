package param

import (
	"context"
	"net/http"
)

// Param ...
type Param string

// SetParam ...
func SetParam(r *http.Request, k, v string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), Param(k), v))
}

// GetParam returns param k value
func GetParam(r *http.Request, k string) string {
	param, _ := r.Context().Value(Param(k)).(string)
	return param
}
