package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/nagaremono/buckler/internals"
)

func handleEcho(c net.Conn, req *internals.Request) {
	arg, found := strings.CutPrefix(req.Target, "/echo/")
	if !found {
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     500,
			StatusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
		return
	}
	body := []byte(arg)

	res := &internals.Response{
		Status:     200,
		StatusText: "OK",
		Protocol:   "HTTP/1.1",
		Headers: []string{
			"Content-Type: text/plain",
			"Content-Length: " + fmt.Sprintf("%d", len(body)),
		},
		Body: body,
	}

	_, err := c.Write(res.Bytes())
	if err != nil {
		fmt.Println(err)
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     500,
			StatusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
	}
}

func handleUserAgent(c net.Conn, req *internals.Request) {
	body := req.Headers["User-Agent"]
	res := &internals.Response{
		Status:     200,
		StatusText: "OK",
		Protocol:   "HTTP/1.1",
		Headers: []string{
			"Content-Type: text/plain",
			"Content-Length: " + fmt.Sprintf("%d", len(body)),
		},
		Body: []byte(body),
	}

	_, err := c.Write([]byte(res.String()))
	if err != nil {
		fmt.Println(err)
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     500,
			StatusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
	}
}

func handleReadFile(c net.Conn, req *internals.Request) {
	filename, found := strings.CutPrefix(req.Target, "/files/")
	if !found {
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     404,
			StatusText: "Not Found",
		}
		c.Write(res.Bytes())
		return
	}
	file, err := os.Open(path.Join(*dirFlag, filename))
	if err != nil {
		if os.IsNotExist(err) {
			res := &internals.Response{
				Protocol:   "HTTP/1.1",
				Status:     404,
				StatusText: "Not Found",
			}
			c.Write(res.Bytes())
			return
		}
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     500,
			StatusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			res := &internals.Response{
				Protocol:   "HTTP/1.1",
				Status:     500,
				StatusText: "Internal Server Error",
			}
			c.Write(res.Bytes())
		}
	}()

	body := []byte{}
	buf := make([]byte, 1024)

	for {
		nRead, err := file.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			res := &internals.Response{
				Protocol:   "HTTP/1.1",
				Status:     500,
				StatusText: "Internal Server Error",
			}
			c.Write(res.Bytes())
			return
		}
		if nRead == 0 {
			break
		}
		body = append(body, buf[:nRead]...)
	}

	res := &internals.Response{
		Status:     200,
		StatusText: "OK",
		Protocol:   "HTTP/1.1",
		Headers: []string{
			"Content-Type: application/octet-stream",
			"Content-Length: " + fmt.Sprintf("%d", len(body)),
		},
		Body: []byte(body),
	}

	_, err = c.Write([]byte(res.String()))
	if err != nil {
		fmt.Println(err)
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     500,
			StatusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
	}
}

func handleWriteFile(c net.Conn, req *internals.Request) {
	filename, found := strings.CutPrefix(req.Target, "/files/")
	if !found {
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     404,
			StatusText: "Not Found",
		}
		c.Write(res.Bytes())
		return
	}
	file, err := os.Create(path.Join(*dirFlag, filename))
	if err != nil {
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     500,
			StatusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			res := &internals.Response{
				Protocol:   "HTTP/1.1",
				Status:     500,
				StatusText: "Internal Server Error",
			}
			c.Write(res.Bytes())
		}
	}()

	bodyLen, _ := strconv.Atoi(req.Headers["Content-Length"])
	_, err = file.Write(req.Body[:bodyLen])
	if err != nil {
		fmt.Println(err)
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     500,
			StatusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
		return
	}

	res := &internals.Response{
		Status:     201,
		StatusText: "Created",
		Protocol:   "HTTP/1.1",
		Headers: []string{
			"Content-Length: 0",
		},
	}

	_, err = c.Write([]byte(res.String()))
	if err != nil {
		fmt.Println(err)
		res := &internals.Response{
			Protocol:   "HTTP/1.1",
			Status:     500,
			StatusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
	}
}
