package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

func handleEcho(c net.Conn, req *Request) {
	arg, found := strings.CutPrefix(req.target, "/echo/")
	if !found {
		res := &Response{
			protocol:   "HTTP/1.1",
			status:     500,
			statusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
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

	_, err := c.Write(res.Bytes())
	if err != nil {
		fmt.Println(err)
		res := &Response{
			protocol:   "HTTP/1.1",
			status:     500,
			statusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
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
		res := &Response{
			protocol:   "HTTP/1.1",
			status:     500,
			statusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
	}
}

func handleFiles(c net.Conn, req *Request) {
	filename, found := strings.CutPrefix(req.target, "/files/")
	if !found {
		res := &Response{
			protocol:   "HTTP/1.1",
			status:     404,
			statusText: "Not Found",
		}
		c.Write(res.Bytes())
		return
	}
	if req.method == "GET" {
		handleReadFile(filename, c)
	} else if req.method == "POST" {
		handleWriteFile(filename, req, c)
	}
}

func handleReadFile(filename string, c net.Conn) {
	file, err := os.Open(path.Join(*dirFlag, filename))
	if err != nil {
		if os.IsNotExist(err) {
			res := &Response{
				protocol:   "HTTP/1.1",
				status:     404,
				statusText: "Not Found",
			}
			c.Write(res.Bytes())
			return
		}
		res := &Response{
			protocol:   "HTTP/1.1",
			status:     500,
			statusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			res := &Response{
				protocol:   "HTTP/1.1",
				status:     500,
				statusText: "Internal Server Error",
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
			res := &Response{
				protocol:   "HTTP/1.1",
				status:     500,
				statusText: "Internal Server Error",
			}
			c.Write(res.Bytes())
			return
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
		res := &Response{
			protocol:   "HTTP/1.1",
			status:     500,
			statusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
	}
}

func handleWriteFile(filename string, req *Request, c net.Conn) {
	file, err := os.Create(path.Join(*dirFlag, filename))
	if err != nil {
		res := &Response{
			protocol:   "HTTP/1.1",
			status:     500,
			statusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			res := &Response{
				protocol:   "HTTP/1.1",
				status:     500,
				statusText: "Internal Server Error",
			}
			c.Write(res.Bytes())
		}
	}()

	bodyLen, _ := strconv.Atoi(req.headers["Content-Length"])
	_, err = file.Write(req.body[:bodyLen])
	if err != nil {
		fmt.Println(err)
		res := &Response{
			protocol:   "HTTP/1.1",
			status:     500,
			statusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
		return
	}

	res := &Response{
		status:     201,
		statusText: "Created",
		protocol:   "HTTP/1.1",
		headers: []string{
			"Content-Length: 0",
		},
	}

	_, err = c.Write([]byte(res.String()))
	if err != nil {
		fmt.Println(err)
		res := &Response{
			protocol:   "HTTP/1.1",
			status:     500,
			statusText: "Internal Server Error",
		}
		c.Write(res.Bytes())
	}
}
