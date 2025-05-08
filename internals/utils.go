package internals

import (
	"net"
	"strings"
)

func ParseRequest(c net.Conn) (*Request, error) {
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
		Method:   requestLine[0],
		Target:   requestLine[1],
		Protocol: requestLine[2],
		Headers:  headers,
		Body:     []byte(body),
	}, nil
}
