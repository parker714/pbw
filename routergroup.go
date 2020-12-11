package pbw

type RouterGroup interface {
	GET(pattern string, handler HandlerFunc)
	POST(pattern string, handler HandlerFunc)
	Use(middleware ...HandlerFunc)
}

type routerGroup struct {
	engine Engine

	prefix     string
	middleware []HandlerFunc
}

func (r *routerGroup) GET(pattern string, handler HandlerFunc) {
	r.engine.GET(r.prefix+pattern, handler)
}

func (r *routerGroup) POST(pattern string, handler HandlerFunc) {
	r.engine.POST(r.prefix+pattern, handler)
}

func (r *routerGroup) Use(middleware ...HandlerFunc) {
	r.middleware = append(r.middleware, middleware...)
}
