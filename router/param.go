package router

import (
	"net/http"

	"devcode.xeemore.com/systech/gojunkyard/router/internal/param"
)

// GetParam returns param k value
func GetParam(r *http.Request, k string) string {
	return param.GetParam(r, k)
}
