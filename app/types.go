package main

import (
	"fmt"
	"net"
	"strings"
)

type Handler func(c net.Conn, r *Request)

type Request struct {
	method   string
	target   string
	protocol string
	headers  map[string]string
	body     []byte
}

func parseRequest(c net.Conn) (*Request, error) {
	data := make([]byte, 512)
	reqString := ""

	for {
		read, err := c.Read(data)
		if err != nil {
			return nil, err
		}
		reqString += string(data)
		if read < 512 {
			break
		}
	}

	requestLineStr, headerAndBody, _ := strings.Cut(reqString, "\r\n")
	requestLine := strings.Split(requestLineStr, " ")

	headerStr, body, _ := strings.Cut(headerAndBody, "\r\n\r\n")
	headers := map[string]string{}

	for _, v := range strings.Split(headerStr, "\r\n") {
		header := strings.Split(v, ": ")
		headers[header[0]] = header[1]
	}

	return &Request{
		method:   requestLine[0],
		target:   requestLine[1],
		protocol: requestLine[2],
		headers:  headers,
		body:     []byte(body),
	}, nil
}

type Response struct {
	status     int
	statusText string
	protocol   string
	headers    []string
	body       []byte
}

func (r *Response) String() string {
	s := fmt.Sprintf("%s %d %s\r\n%s\r\n\r\n%s\r\n\r\n", r.protocol, r.status, r.statusText, strings.Join(r.headers, "\r\n"), string(r.body))
	return s
}

func (r *Response) Bytes() []byte {
	return []byte(r.String())
}

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

func (r *Router) register(method string, target string, handler Handler) {
	r.rHandlers = append(r.rHandlers, &RouteHandler{
		target,
		method,
		handler,
	})
}

func (r *Router) match(req *Request, c net.Conn) {
	var handler Handler
	for _, rh := range r.rHandlers {
		if matcher(req, rh) {
			handler = rh.handler
		}
	}

	if handler == nil {
		res := &Response{
			status:     404,
			statusText: "Not Found",
			protocol:   "HTTP/1.1",
		}

		_, err := c.Write(res.Bytes())
		if err != nil {
			fmt.Println(err)
			res := &Response{
				status:     500,
				statusText: "Internal Server Error",
				protocol:   "HTTP/1.1",
			}
			c.Write(res.Bytes())
		}
	}

	handler(c, req)
}

func matcher(req *Request, rh *RouteHandler) bool {
	if req.method != rh.method {
		return false
	}
	targetParts := strings.Split(req.target, "/")
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
