package main

import (
	"fmt"
	"net"
	"strings"
)

type Handler func(c net.Conn)

type Request struct {
	method   string
	target   string
	protocol string
	headers  map[string]string
	body     []byte
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
