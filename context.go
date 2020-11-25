package pbw

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Path       string
	Method     string
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) SetStatusCode(statusCode int) {
	c.StatusCode = statusCode
	c.Writer.WriteHeader(statusCode)
}

func (c *Context) String(code int, value string) {
	c.SetStatusCode(code)
	c.SetHeader("Content-Type", "text/plain")
	_, _ = c.Writer.Write([]byte(value))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetStatusCode(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetStatusCode(code)
	c.SetHeader("Content-Type", "text/html")
	_, _ = c.Writer.Write([]byte(html))
}

func (c *Context) Data(code int, data []byte) {
	c.SetStatusCode(code)
	_, _ = c.Writer.Write(data)
}
