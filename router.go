package pbw

import (
	"net/http"
	"strings"
)

const (
	patternPrefixErr = "router: addPattern pattern must begin with /"
	httpNotFound     = "router: 404 NotFound"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addPattern(method string, pattern string, handler HandlerFunc) {
	if !strings.HasPrefix(pattern, "/") {
		panic(patternPrefixErr)
	}

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = newNode("")
	}

	patterns := r.parsePattern(pattern)
	r.roots[method].insert(patterns)

	r.handlers[method+"-"+strings.Join(patterns, "/")] = handler
}

func (r *router) parsePattern(pattern string) []string {
	patterns := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, pattern := range patterns {
		parts = append(parts, pattern)
		if strings.HasPrefix(pattern, "*") {
			break
		}
	}
	return parts
}

func (r *router) getPattern(method string, requestPath string) (patterns []string) {
	if _, ok := r.roots[method]; !ok {
		return make([]string, 0)
	}

	return r.roots[method].search(strings.Split(requestPath, "/"), 0)
}

func (r *router) handle(c Context) {
	patterns := r.getPattern(c.Method(), c.Path())
	if 0 == len(patterns) {
		c.Data(http.StatusNotFound, []byte(httpNotFound))
		return
	}

	// add context params
	requestParts := strings.Split(c.Path(), "/")
	for index, part := range patterns {
		if "" == part {
			continue
		}
		if part[0] == ':' {
			c.SetParam(part[1:], requestParts[index])
		}
		if part[0] == '*' && len(part) > 1 {
			c.SetParam(part[1:], strings.Join(requestParts[index:], "/"))
			break
		}
	}

	// router handler
	key := c.Method() + "-" + strings.Join(patterns, "/")
	c.AddHandlers(r.handlers[key])

	// next handler
	c.Next()
}
