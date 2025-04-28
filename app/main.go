package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handle(conn)
	}
}

var dirFlag = flag.String("directory", "", "")

func handle(c net.Conn) {
	req, err := parseRequest(c)
	if err != nil {
		fmt.Println(err)
		res := &Response{
			status:     500,
			statusText: "Internal Server Error",
			protocol:   "HTTP/1.1",
		}
		c.Write(res.Bytes())
		return
	}

	var res *Response
	if req.target == "/" {
		res = &Response{
			status:     200,
			statusText: "OK",
			protocol:   "HTTP/1.1",
		}
	} else if strings.HasPrefix(req.target, "/echo/") {
		handleEcho(c, req)
		return
	} else if strings.HasPrefix(req.target, "/user-agent") {
		handleUserAgent(c, req)
		return
	} else if strings.HasPrefix(req.target, "/files") {
		handleFiles(c, req)
		return
	} else {
		res = &Response{
			status:     404,
			statusText: "Not Found",
			protocol:   "HTTP/1.1",
		}
	}

	_, err = c.Write(res.Bytes())
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
