package internals

import (
	"net"
	"strings"
)

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

	for v := range strings.SplitSeq(headerStr, "\r\n") {
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
