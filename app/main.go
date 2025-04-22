package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
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
		go func() {
			handle(conn)
			conn.Close()
		}()
	}
}

func handle(c net.Conn) {
	req, err := parseRequest(c)
	if err != nil {
		fmt.Println(err)
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
		return
	}

	var res string
	if req.target == "/" {
		res = "HTTP/1.1 200 OK\r\n\r\n"
	} else if strings.HasPrefix(req.target, "/echo/") {
		handleEcho(c, req)
		return
	} else {
		res = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = c.Write([]byte(res))
	if err != nil {
		fmt.Println(err)
		c.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
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

	r := strings.Split(reqString, "\r\n")
	if len(r) < 3 {
		return nil, errors.New("invalid request")
	}

	requestLine := strings.Split(r[0], " ")

	return &Request{
		method:   requestLine[0],
		target:   requestLine[1],
		protocol: requestLine[2],
	}, nil
}
