package main

import (
	"fmt"
	"net"
	"strings"
)

func handleEcho(c net.Conn, req *Request) {
	arg, found := strings.CutPrefix(req.target, "/echo/")
	if !found {
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}
	body := []byte(arg)

	res := &Response{
		status:     200,
		statusText: "OK",
		protocol:   "HTTP/1.1",
		headers: []string{
			"Content-Type: text/plain",
			"Content-Length: " + fmt.Sprintf("%d", len(body)),
		},
		body: body,
	}

	_, err := c.Write([]byte(res.String()))
	if err != nil {
		fmt.Println(err)
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}
}
