package valkyrie

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HTTPRouter holds all core method which have projectID validation method
type HTTPRouter struct {
	*core
}

// NewHTTPRouter returns projectID validator for juliendschmidth httprouter
func NewHTTPRouter(option *Option) *HTTPRouter {
	return &HTTPRouter{construct(option)}
}

// AuthProject is the middleware of httprouter used to validate projectID
func (c *HTTPRouter) AuthProject(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		pid, ok := c.authProject(r)
		if ok {
			r = r.WithContext(context.WithValue(r.Context(), PID, pid))
			h(w, r, ps)
			return
		}
		pid, ok = c.authHostname(r)
		if ok {
			r = r.WithContext(context.WithValue(r.Context(), PID, pid))
			h(w, r, ps)
			return
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
