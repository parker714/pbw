package pbw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Handle(t *testing.T) {
	t.Run("method not exit", func(t *testing.T) {
		engine := New()
		engine.GET("/user", nil)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/user", nil)
		engine.ServeHTTP(response, request)

		if http.StatusNotFound != response.Code {
			t.Fatalf("expect %d, got %d", http.StatusNotFound, response.Code)
		}
		if NotFound != response.Body.String() {
			t.Fatalf("expect %s, got %s", NotFound, response.Body.String())
		}
	})

	t.Run("path not exit", func(t *testing.T) {
		engine := New()
		engine.GET("/user", nil)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/abc", nil)
		engine.ServeHTTP(response, request)

		if NotFound != response.Body.String() {
			t.Fatalf("expect %s, got %s", NotFound, response.Body.String())
		}
	})

	t.Run("context param", func(t *testing.T) {
		engine := New()
		engine.GET("/user/:name/*action", func(c Context) {
			if c.Param("name") != "pb" {
				t.Fatalf("expect %s, got %s", "pb", c.Param("name"))
			}
			c.Data(http.StatusOK, []byte("pb"))
		})

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/user/pb/say", nil)
		engine.ServeHTTP(response, request)
		if http.StatusOK != response.Code {
			t.Fatalf("expect %d, got %d", http.StatusOK, response.Code)
		}
		if "pb" != response.Body.String() {
			t.Fatalf("expect %s, got %s", "pb", response.Body.String())
		}
	})
}
