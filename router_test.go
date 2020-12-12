package pbw

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRouter_addPattern(t *testing.T) {
	t.Run("params invalid", func(t *testing.T) {
		defer func() {
			err := recover()
			if patternPrefixErr != err {
				t.Fatalf("expect %v, got %v\n", patternPrefixErr, err)
			}
		}()
		mr := newRouter()
		mr.addPattern("", "", nil)
	})

	t.Run("params normal", func(t *testing.T) {
		r := newRouter()
		r.addPattern("GET", "/", nil)
		expect := &router{
			roots: map[string]*node{
				"GET": {
					patterns: make([]string, 0),
					next: map[string]*node{
						"": {
							patterns: make([]string, 0),
							next: map[string]*node{
								"": {
									patterns: []string{"", ""},
									next:     make(map[string]*node),
								},
							},
						},
					},
				},
			},
			handlers: map[string]HandlerFunc{
				"GET-/": nil,
			},
		}
		if !reflect.DeepEqual(expect, r) {
			t.Fatalf("/: expect %v, got %v\n", expect, r)
		}

		r.addPattern("GET", "/user", nil)
		expect.roots["GET"].next[""].next["user"] = &node{
			patterns: []string{"", "user"},
			next:     make(map[string]*node),
		}
		expect.handlers["GET-/user"] = nil
		if !reflect.DeepEqual(expect, r) {
			t.Fatalf("/user: expect %v, got %v\n", expect, r)
		}

		r.addPattern("GET", "/:lang/doc", nil)
		expect.roots["GET"].next[""].next[":lang"] = &node{
			isWild:   true,
			patterns: make([]string, 0),
			next: map[string]*node{
				"doc": {
					patterns: []string{"", ":lang", "doc"},
					next:     make(map[string]*node),
				},
			},
		}
		expect.handlers["GET-/:lang/doc"] = nil
		if !reflect.DeepEqual(expect, r) {
			t.Fatalf("/:lang/doc: expect %v, got %v\n", expect, r)
		}
	})
}

func TestRouter_parsePattern(t *testing.T) {

}

func TestRouter_getPattern(t *testing.T) {
	t.Run("not exist method", func(t *testing.T) {
		mr := newRouter()
		mr.addPattern("GET", "/", nil)

		cases := []struct {
			method      string
			requestPath string
			except      []string
		}{
			{
				method:      "GET",
				requestPath: "/",
				except:      []string{"", ""},
			},
			{
				method:      "PB",
				requestPath: "/",
				except:      make([]string, 0),
			},
		}

		for _, c := range cases {
			if !reflect.DeepEqual(c.except, mr.getPattern(c.method, c.requestPath)) {
				t.Fatalf("expect %v, got %v\n", c.except, mr.getPattern(c.method, c.requestPath))
			}
		}
	})

	t.Run("not exist router", func(t *testing.T) {
		mr := newRouter()
		mr.addPattern("GET", "/", nil)

		got := mr.getPattern("GET", "/user")
		expect := make([]string, 0)
		if !reflect.DeepEqual(expect, got) {
			t.Fatalf("expect %v, got %v\n", expect, got)
		}
	})

	t.Run("router normal", func(t *testing.T) {
		mr := newRouter()
		mr.addPattern("GET", "/", nil)
		mr.addPattern("GET", "/user", nil)
		mr.addPattern("GET", "/:lang/doc", nil)

		cases := []struct {
			method      string
			requestPath string
			except      []string
		}{
			{
				method:      "GET",
				requestPath: "/user",
				except:      []string{"", "user"},
			},
			{
				method:      "GET",
				requestPath: "/golang/doc",
				except:      []string{"", ":lang", "doc"},
			},
		}

		for _, c := range cases {
			if !reflect.DeepEqual(c.except, mr.getPattern(c.method, c.requestPath)) {
				t.Fatalf("expect %v, got %v\n", c.except, mr.getPattern(c.method, c.requestPath))
			}
		}
	})
}

func TestRouter_handle(t *testing.T) {
	t.Run("router no exit", func(t *testing.T) {
		mr := newRouter()
		mr.addPattern("GET", "/", nil)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/user", nil)
		c := NewContext(response, request)
		mr.handle(c)

		if response.Body.String() != httpNotFound {
			t.Fatalf("expect %s, got %s\n", httpNotFound, response.Body.String())
		}
	})

	t.Run("router handle", func(t *testing.T) {
		mr := newRouter()
		mr.addPattern("GET", "/:lang/doc/*name", func(c Context) {})

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/go/doc/pb", nil)
		c := NewContext(response, request, func(c Context) {})
		mr.handle(c)

		if "go" != c.Param("lang") {
			t.Fatalf("expect %s, got %s\n", "go", c.Param("lang"))
		}
		if "pb" != c.Param("name") {
			t.Fatalf("expect %s, got %s\n", "pb", c.Param("name"))
		}
	})
}
