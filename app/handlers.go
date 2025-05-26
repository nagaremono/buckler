package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/nagaremono/buckler/internals"
)

func handleEcho(req *internals.Request, res *internals.Response) {
	arg, found := strings.CutPrefix(req.Target, "/echo/")
	if !found {
		res.Status = 500
		res.StatusText = "Internal Server Error"
		return
	}
	body := []byte(arg)

	res.Status = 200
	res.StatusText = "OK"
	res.Headers = append(res.Headers,
		[]string{
			"Content-Type: text/plain",
			"Content-Length: " + fmt.Sprintf("%d", len(body)),
		}...)
	res.Body = body
}

func handleUserAgent(req *internals.Request, res *internals.Response) {
	body := req.Headers["User-Agent"]
	res.Status = 200
	res.StatusText = "OK"
	res.Headers = []string{
		"Content-Type: text/plain",
		"Content-Length: " + fmt.Sprintf("%d", len(body)),
	}
	res.Body = []byte(body)
}

func handleReadFile(req *internals.Request, res *internals.Response) {
	filename, found := strings.CutPrefix(req.Target, "/files/")
	if !found {
		res.Status = 404
		res.StatusText = "Not Found"
		return
	}
	file, err := os.Open(path.Join(*dirFlag, filename))
	if err != nil {
		if os.IsNotExist(err) {
			res.Status = 404
			res.StatusText = "Not Found"
		} else {
			res.Status = 500
			res.StatusText = "Internal Server Error"
		}
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			res.Status = 500
			res.StatusText = "Internal Server Error"
		}
	}()

	body := []byte{}
	buf := make([]byte, 1024)

	for {
		nRead, err := file.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			res.Status = 500
			res.StatusText = "Internal Server Error"
			return
		}
		if nRead == 0 {
			break
		}
		body = append(body, buf[:nRead]...)
	}

	res.Status = 200
	res.StatusText = "OK"
	res.Headers = append(res.Headers, []string{
		"Content-Type: application/octet-stream",
		"Content-Length: " + fmt.Sprintf("%d", len(body)),
	}...)
	res.Body = []byte(body)
}

func handleWriteFile(req *internals.Request, res *internals.Response) {
	filename, found := strings.CutPrefix(req.Target, "/files/")
	if !found {
		res.Status = 404
		res.StatusText = "Not Found"
		return
	}
	file, err := os.Create(path.Join(*dirFlag, filename))
	if err != nil {
		res.Status = 500
		res.StatusText = "Internal Server Error"
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			res.Status = 500
			res.StatusText = "Internal Server Error"
		}
	}()

	bodyLen, _ := strconv.Atoi(req.Headers["Content-Length"])
	_, err = file.Write(req.Body[:bodyLen])
	if err != nil {
		fmt.Println(err)
		res.Status = 500
		res.StatusText = "Internal Server Error"
		return
	}

	res.Status = 201
	res.StatusText = "Created"
	res.Headers = append(res.Headers, "Content-Length: 0")
}
