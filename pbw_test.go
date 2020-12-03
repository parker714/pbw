package pbw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEngine_Run(t *testing.T) {
	engine := New()
	err := engine.Run("80000")
	if err == nil {
		t.Fatalf("expect nil, got %s", err)
	}
}

func TestEngine_GET(t *testing.T) {
	engine := New()
	engine.GET("/user", func(c Context) {
		c.Data(http.StatusOK, []byte("pb"))
	})
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user", nil)
	engine.ServeHTTP(response, request)

	if "pb" != response.Body.String() {
		t.Fatalf("expect %s, got %s", "pb", response.Body.String())
	}
}

func TestEngine_POST(t *testing.T) {
	engine := New()
	engine.POST("/user", func(c Context) {
		c.Data(http.StatusOK, []byte("pb"))
	})
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/user", nil)
	engine.ServeHTTP(response, request)

	if "pb" != response.Body.String() {
		t.Fatalf("expect %s, got %s", "pb", response.Body.String())
	}
}

func TestContext_SetHeader(t *testing.T) {
	engine := New()
	engine.GET("/user", func(c Context) {
		c.SetHeader("token", "pb-token")
		c.Data(http.StatusOK, []byte("pb"))
	})
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user", nil)
	engine.ServeHTTP(response, request)

	if "pb-token" != response.Header().Get("token") {
		t.Fatalf("expect %s, got %s", "pb-token", response.Header().Get("token"))
	}
	if "pb" != response.Body.String() {
		t.Fatalf("expect %s, got %s", "pb", response.Body.String())
	}
}

func TestContext_JSON(t *testing.T) {
	engine := New()
	engine.GET("/user", func(c Context) {
		c.JSON(http.StatusOK, H{"name": "pb"})
	})
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user", nil)
	engine.ServeHTTP(response, request)

	if response.Body.String() != `{"name":"pb"}` {
		t.Fatalf("expect %s, got %s", `{"name":"pb"}`, response.Body.String())
	}
}

func TestContext_JSON_Panic(t *testing.T) {
	t.Run("json marshal err", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Fatalf("expect %s, got nil", "exception")
			}
		}()

		engine := New()
		engine.GET("/user", func(c Context) {
			c.JSON(http.StatusOK, make(chan int))
		})
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/user", nil)
		engine.ServeHTTP(response, request)
	})
}
