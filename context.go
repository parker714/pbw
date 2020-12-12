package pbw

import (
	"encoding/json"
	"net/http"
)

// H http header params.
type H map[string]interface{}

// Context used to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response.
type Context interface {
	Method() string
	Path() string
	Query(key string) string
	SetParam(key, value string)
	Param(key string) string
	PostForm(key string) string
	SetStatusCode(statusCode int)
	SetHeader(key string, value string)
	Data(code int, data []byte)
	JSON(code int, obj interface{})
	AddHandlers(handlers ...HandlerFunc)
	Next()
}

type context struct {
	writer http.ResponseWriter
	req    *http.Request

	// request uri params
	params map[string]string

	// middleware、handler callback
	index    int
	handlers []HandlerFunc
}

// NewContext used return context.
func NewContext(w http.ResponseWriter, req *http.Request, hfs ...HandlerFunc) Context {
	return &context{
		writer:   w,
		req:      req,
		params:   make(map[string]string),
		handlers: hfs,
		index:    -1,
	}
}

func (c *context) Method() string {
	return c.req.Method
}

func (c *context) Path() string {
	return c.req.URL.Path
}

func (c *context) Query(key string) string {
	return c.req.URL.Query().Get(key)
}

func (c *context) SetParam(key, value string) {
	c.params[key] = value
}

func (c *context) Param(key string) string {
	value, _ := c.params[key]
	return value
}

func (c *context) PostForm(key string) string {
	return c.req.FormValue(key)
}

func (c *context) SetStatusCode(statusCode int) {
	c.writer.WriteHeader(statusCode)
}

func (c *context) SetHeader(key string, value string) {
	c.writer.Header().Set(key, value)
}

func (c *context) Data(code int, data []byte) {
	c.SetStatusCode(code)
	_, _ = c.writer.Write(data)
}

func (c *context) JSON(code int, obj interface{}) {
	c.SetStatusCode(code)
	c.SetHeader("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	if _, err = c.writer.Write(jsonBytes); err != nil {
		panic(err)
	}
}

func (c *context) AddHandlers(handlers ...HandlerFunc) {
	c.handlers = append(c.handlers, handlers...)
}

func (c *context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}
