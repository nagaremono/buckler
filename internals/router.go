package internals

import (
	"fmt"
	"net"
	"strings"
)

type Handler func(c net.Conn, r *Request, s *Response)

type RouteHandler struct {
	target  string
	method  string
	handler Handler
}

type Router struct {
	rHandlers []*RouteHandler
}

func NewRouter() *Router {
	return &Router{
		rHandlers: []*RouteHandler{},
	}
}

func (r *Router) Register(method string, target string, handler Handler) {
	r.rHandlers = append(r.rHandlers, &RouteHandler{
		target,
		method,
		handler,
	})
}

func (r *Router) Match(req *Request, c net.Conn) Handler {
	var handler Handler
	for _, rh := range r.rHandlers {
		if Matcher(req, rh) {
			handler = rh.handler
		}
	}

	res := &Response{
		Protocol: "HTTP/1.1",
	}

	if handler == nil {
		res.Status = 404
		res.StatusText = "Not Found"

		_, err := c.Write(res.Bytes())
		if err != nil {
			fmt.Println(err)
			res.Status = 500
			res.StatusText = "Internal Server Error"
			res.Protocol = "HTTP/1.1"
		}
		c.Write(res.Bytes())
		return nil
	}

	return handler
}

func Matcher(req *Request, rh *RouteHandler) bool {
	if req.Method != rh.method {
		return false
	}
	targetParts := strings.Split(req.Target, "/")
	handlerParts := strings.Split(rh.target, "/")
	for i := 0; i < len(targetParts); i++ {
		if i >= len(handlerParts) {
			return false
		}
		if !strings.HasPrefix(handlerParts[i], ":") && targetParts[i] != handlerParts[i] {
			return false
		}
	}

	return true
}
