package pbw

import (
	"net/http"
	"strings"
)

const NotFound = "404 NOT FOUND"

type Router interface {
	AddRouter(method string, pattern string, handler HandlerFunc)
	Handle(c Context)
}

type router struct {
	roots    map[string]Node
	handlers map[string]HandlerFunc
}

func NewRouter() Router {
	return &router{
		roots: map[string]Node{
			http.MethodGet:  NewNode(),
			http.MethodPost: NewNode(),
		},
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) AddRouter(method string, pattern string, handler HandlerFunc) {
	parts := r.parsePattern(pattern)
	r.roots[method].Insert(parts)

	key := method + "-" + strings.Join(parts, "/")
	r.handlers[key] = handler
}

// parsePattern used to resolve registered routes
// /welcome            -> welcome
// /user/:name         -> user/:name
// /user/:name/*action -> user/:name/*action
func (r *router) parsePattern(pattern string) []string {
	patterns := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range patterns {
		if item == "" {
			continue
		}
		parts = append(parts, item)
		if strings.HasPrefix(item, "*") {
			break
		}
	}
	return parts
}

func (r *router) Handle(c Context) {
	root, ok := r.roots[c.Method()]
	if !ok {
		c.Data(http.StatusNotFound, []byte(NotFound))
		return
	}

	pattern := root.Search(c.Path())
	if pattern == "" {
		c.Data(http.StatusNotFound, []byte(NotFound))
		return
	}

	// add context params
	patternParts := strings.Split(pattern, "/")
	requestParts := strings.Split(c.Path(), "/")[1:]
	for index, part := range patternParts {
		if part[0] == ':' {
			c.SetParam(part[1:], requestParts[index])
		}
		if part[0] == '*' && len(part) > 1 {
			c.SetParam(part[1:], strings.Join(requestParts[index:], "/"))
			break
		}
	}

	// router handler
	key := c.Method() + "-" + pattern
	r.handlers[key](c)
}
