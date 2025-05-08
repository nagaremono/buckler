package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/nagaremono/buckler/internals"
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

var router *internals.Router

func initRouter() {
	router = internals.NewRouter()
	router.Register("GET", "/echo/:s", handleEcho)
	router.Register("GET", "/user-agent", handleUserAgent)
	router.Register("GET", "/files/:file", handleReadFile)
	router.Register("POST", "/files/:file", handleWriteFile)
}

func handle(c net.Conn) {
	req, err := internals.ParseRequest(c)
	if err != nil {
		fmt.Println(err)
		res := &internals.Response{
			Status:     500,
			StatusText: "Internal Server Error",
			Protocol:   "HTTP/1.1",
		}
		c.Write(res.Bytes())
		return
	}

	var res *internals.Response
	if req.Target == "/" {
		res = &internals.Response{
			Status:     200,
			StatusText: "OK",
			Protocol:   "HTTP/1.1",
		}
		c.Write(res.Bytes())
		return
	}

	router.Match(req, c)
}
