package pbw

import (
	"net/http"
	"strings"
)

type HandlerFunc func(Context)

type Engine interface {
	Use(middleware ...HandlerFunc)
	Group(prefix string) RouterGroup
	GET(pattern string, handler HandlerFunc)
	POST(pattern string, handler HandlerFunc)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	Run(addr string) (err error)
}

type engine struct {
	groups []*routerGroup
	router *router
}

func New() Engine {
	return &engine{
		router: newRouter(),
		groups: make([]*routerGroup, 0),
	}
}

func (e *engine) Use(middleware ...HandlerFunc) {
	e.Group("/").Use(middleware...)
}

func (e *engine) Group(prefix string) RouterGroup {
	routerGroup := &routerGroup{
		engine: e,
		prefix: prefix,
	}
	e.groups = append(e.groups, routerGroup)
	return routerGroup
}

func (e *engine) GET(pattern string, handler HandlerFunc) {
	e.router.addPattern(http.MethodGet, pattern, handler)
}

func (e *engine) POST(pattern string, handler HandlerFunc) {
	e.router.addPattern(http.MethodPost, pattern, handler)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middleware []HandlerFunc

	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middleware = append(middleware, group.middleware...)
		}
	}
	e.router.handle(NewContext(w, req, middleware...))
}

func (e *engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
