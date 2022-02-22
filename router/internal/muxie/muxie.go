package muxie

import (
	"context"
	"net/http"

	"devcode.xeemore.com/systech/gojunkyard/router/internal/param"

	"github.com/kataras/muxie"
)

// Router is engine which composed from muxie
// This engine is recomended to use if there is conflicted routes which cannot handled by httprouter
// Pros:
// * It can handle same pattern routes. The priority is depends on order definition
// Cons:
// * Not many review for this router. But the developer claims that the performance better than httprouter
type Router struct {
	*muxie.Mux
}

// New returns muxie engine
func New() *Router {
	return &Router{muxie.NewMux()}
}

// Handle handles http request by method which is specified by the user
func (r *Router) Handle(method, path string, handler http.Handler) {
	var methodHandler = Methods()

	// if path already registered, then use older method handler
	if node := r.Mux.Routes.Search(path, new(nopParamSetter)); node != nil {
		methodHandler = node.Handler.(*MethodHandler)
	}

	r.Mux.Handle(path, methodHandler.HandleFunc(method, func(w http.ResponseWriter, req *http.Request) {
		var ctx = req.Context()
		for _, p := range muxie.GetParams(w) {
			ctx = context.WithValue(ctx, param.Param(p.Key), p.Value)
		}
		handler.ServeHTTP(w, req.WithContext(ctx))
	}))
}

// Lookup check if path already registered
func (r *Router) Lookup(method, path string) bool {
	node := r.Mux.Routes.Search(path, new(nopParamSetter))
	if node == nil {
		return false
	}
	return node.Handler.(*MethodHandler).Lookup(method)
}
