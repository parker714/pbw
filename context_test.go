package pbw

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mc Context

func init() {
	mc = mockContext()
}

func mockContext() Context {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	_ = w.WriteField("age", "18")
	_ = w.Close()
	request := httptest.NewRequest(http.MethodGet, "/user?name=pb", body)
	request.Header.Set("Content-Type", w.FormDataContentType())
	return NewContext(nil, request)
}

func TestContext_Method(t *testing.T) {
	if http.MethodGet != mc.Method() {
		t.Fatalf("expect %s, got %s", http.MethodGet, mc.Method())
	}
}

func TestContext_Path(t *testing.T) {
	if "/user" != mc.Path() {
		t.Fatalf("expect %s, got %s", "/user", mc.Path())
	}
}

func TestContext_Query(t *testing.T) {
	if "pb" != mc.Query("name") {
		t.Fatalf("expect %s, got %s", "pb", mc.Query("name"))
	}
}

func TestContext_PostForm(t *testing.T) {
	if "18" != mc.PostForm("age") {
		t.Fatalf("expect %s, got %s", "pb", mc.PostForm("age"))
	}
}
