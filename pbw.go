package pbw

import (
	"net/http"
)

type HandlerFunc func(Context)

type Engine struct {
	router Router
}

func New() *Engine {
	return &Engine{router: NewRouter()}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.Handle(NewContext(w, req))
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.AddRouter(http.MethodGet, pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.AddRouter(http.MethodPost, pattern, handler)
}
