package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"
)

func handleEcho(c net.Conn, req *Request) {
	arg, found := strings.CutPrefix(req.target, "/echo/")
	if !found {
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		return
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

func handleUserAgent(c net.Conn, req *Request) {
	body := req.headers["User-Agent"]
	res := &Response{
		status:     200,
		statusText: "OK",
		protocol:   "HTTP/1.1",
		headers: []string{
			"Content-Type: text/plain",
			"Content-Length: " + fmt.Sprintf("%d", len(body)),
		},
		body: []byte(body),
	}

	_, err := c.Write([]byte(res.String()))
	if err != nil {
		fmt.Println(err)
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}
}

func handleFiles(c net.Conn, req *Request) {
	filename, found := strings.CutPrefix(req.target, "/files/")
	if !found {
		c.Write([]byte("HTTP/1.1 400 Not Found\r\n\r\n"))
		return
	}
	file, err := os.Open(path.Join(*dirFlag, filename))
	if err != nil {
		if os.IsNotExist(err) {
			c.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			return
		}
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		}
	}()

	body := []byte{}
	buf := make([]byte, 1024)

	for {
		nRead, err := file.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		}
		if nRead == 0 {
			break
		}
		body = append(body, buf[:nRead]...)
	}

	res := &Response{
		status:     200,
		statusText: "OK",
		protocol:   "HTTP/1.1",
		headers: []string{
			"Content-Type: application/octet-stream",
			"Content-Length: " + fmt.Sprintf("%d", len(body)),
		},
		body: []byte(body),
	}

	_, err = c.Write([]byte(res.String()))
	if err != nil {
		fmt.Println(err)
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
	}
}
