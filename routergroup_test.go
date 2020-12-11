package pbw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterGroup_GET(t *testing.T) {
	engine := New()
	g1 := engine.Group("/user")
	{
		g1.GET("/order", func(c Context) {
			c.Data(http.StatusOK, []byte("order"))
		})
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user/order", nil)
	engine.ServeHTTP(response, request)
	if "order" != response.Body.String() {
		t.Fatalf("expect %s, got %s", "order", response.Body.String())
	}
}

func TestRouterGroup_POST(t *testing.T) {
	engine := New()
	g1 := engine.Group("/user")
	{
		g1.POST("/buy", func(c Context) {
			c.Data(http.StatusOK, []byte("order"))
		})
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/user/buy", nil)
	engine.ServeHTTP(response, request)
	if "order" != response.Body.String() {
		t.Fatalf("expect %s, got %s", "order", response.Body.String())
	}
}
