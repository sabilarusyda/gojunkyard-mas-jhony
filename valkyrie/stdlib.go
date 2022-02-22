package valkyrie

import (
	"context"
	"net/http"
)

// Stdlib holds all core method which have projectID validation method
type Stdlib struct {
	*core
}

// NewStdlib returns projectID validator for net/http handlefunc
func NewStdlib(option *Option) *Stdlib {
	return &Stdlib{
		construct(option),
	}
}

// AuthProject is the middleware of net/http handlefunc used to validate projectID
func (c *Stdlib) AuthProject(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pid, ok := c.authProject(r)
		if ok {
			r = r.WithContext(context.WithValue(r.Context(), PID, pid))
			h.ServeHTTP(w, r)
			return
		}
		pid, ok = c.authHostname(r)
		if ok {
			r = r.WithContext(context.WithValue(r.Context(), PID, pid))
			h.ServeHTTP(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
