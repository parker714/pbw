package pbw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	engine := New()
	engine.Use(Recovery())
	engine.GET("/user", func(c Context) {
		panic("error")
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user", nil)
	engine.ServeHTTP(response, request)
	if internalServerErr != response.Body.String() {
		t.Fatalf("expect %s, got %s", internalServerErr, response.Body.String())
	}
}
