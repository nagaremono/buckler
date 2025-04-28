package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	initRouter()

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

var router *Router

func initRouter() {
	router = NewRouter()
	router.register("GET", "/echo/:s", handleEcho)
	router.register("GET", "/user-agent", handleUserAgent)
	router.register("GET", "/files/:file", handleFiles)
}

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
		c.Write(res.Bytes())
		return
	}

	router.match(req, c)
}
