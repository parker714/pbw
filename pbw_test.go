package pbw

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngine_ServeHTTP(t *testing.T) {
	t.Run("HTTP 200 router", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/user", nil)
		response := httptest.NewRecorder()

		engine := New()
		engine.GET("/user", func(c *Context) {
			_, _ = fmt.Fprint(c.Writer, "pb")
		})
		engine.POST("/user", func(c *Context) {})
		engine.ServeHTTP(response, request)

		if http.StatusOK != response.Code {
			t.Fatalf("http code err, got %d, want %d", response.Code, http.StatusOK)
		}
		if "pb" != response.Body.String() {
			t.Fatalf("http body err, got %s, want %s", response.Body.String(), "pb")
		}
	})

	t.Run("HTTP 404 router", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/user1", nil)
		response := httptest.NewRecorder()

		engine := New()
		engine.GET("/user", func(c *Context) {})
		engine.ServeHTTP(response, request)

		want := "404 NOT FOUND: /user1\n"
		if want != response.Body.String() {
			t.Fatalf("http body err, got %s, want %s", response.Body.String(), want)
		}
	})
}

func TestEngine_Run(t *testing.T) {
	engine := New()
	err := engine.Run("80000")
	if err == nil {
		t.Fatalf("engine run err, got %s, want not nil", err)
	}
}
