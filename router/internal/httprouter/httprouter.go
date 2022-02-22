package httprouter

import (
	"context"
	"net/http"

	"devcode.xeemore.com/systech/gojunkyard/router/internal/param"

	"github.com/julienschmidt/httprouter"
)

// Router is engine which composed from httprouter
// This engine is recommended by many developer due to stability and performance
// Pros:
// * Performance and Stability
// Cons:
// * This router is using tree algorithm for routing definition. So it will be panic if there are same pattern routes
//   Example: "/users/:id" and "/users/blocked"
type Router struct {
	*httprouter.Router
}

// New returns httprouter engine
func New() *Router {
	return &Router{httprouter.New()}
}

// Handle handles http request by method which is specified by the user
func (r *Router) Handle(method, path string, handler http.Handler) {
	r.Router.Handle(method, path, func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		var ctx = req.Context()
		for _, p := range ps {
			ctx = context.WithValue(ctx, param.Param(p.Key), p.Value)
		}
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}

// Lookup check if path already registered
func (r *Router) Lookup(method, path string) bool {
	h, _, _ := r.Router.Lookup(method, path)
	return h != nil
}
