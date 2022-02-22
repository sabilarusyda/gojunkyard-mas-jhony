package router

import (
	"net/http"
)

type middleware func(next http.Handler) http.Handler

type engine interface {
	Handle(method, path string, handler http.Handler)
	Lookup(method, path string) bool
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// Router ...
type Router struct {
	middlewares []middleware
	path        string
	engine      engine
}

// New returns a new initialized Router.
func New() *Router {
	return NewWithEngine(HTTPRouter)
}

// NewWithEngine ...
func NewWithEngine(et EngineType) *Router {
	return &Router{
		middlewares: make([]middleware, 0),
		engine:      getRouterEngine(et),
	}
}

// Use appends new middleware to current Router.
func (r *Router) Use(m ...middleware) *Router {
	router := &Router{
		engine:      r.engine,
		middlewares: make([]middleware, 0, len(r.middlewares)+len(m)),
		path:        r.path,
	}
	router.middlewares = append(router.middlewares, r.middlewares...)
	router.middlewares = append(router.middlewares, m...)
	return router
}

// Group returns new *Router with given path and middlewares.
// It should be used for handles which have same path prefix or common middlewares.
func (r *Router) Group(path string, m ...middleware) *Router {
	router := r.Use(m...)
	router.path += path
	return router
}

// HandleFunc registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) HandleFunc(method, path string, handler http.HandlerFunc) {
	r.Handle(method, path, handler)
}

// Handle is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *Router) Handle(method, path string, handler http.Handler) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	path = r.path + path
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	r.engine.Handle(method, path, handler)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodHead, path, handler)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodOptions, path, handler)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPatch, path, handler)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.engine.ServeHTTP(w, req)
}

// Static serves files from given root directory.
func (r *Router) Static(path, root string) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path should end with '/*filepath' in path '" + path + "'.")
	}

	var (
		base       = r.path + path[:len(path)-9]
		fileServer = http.StripPrefix(base, http.FileServer(http.Dir(root)))
	)

	r.Handle("GET", path, fileServer)
}
