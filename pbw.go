package pbw

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.handle(newContext(w, req))
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRouter("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRouter("POST", pattern, handler)
}
