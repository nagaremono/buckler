package internals

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"slices"
	"strings"
)

var supportedCompression = []string{"gzip"}

type HeaderHandler func(c net.Conn, r *Request, s *Response)

func CompressionHandler(c net.Conn, r *Request, s *Response) {
	headers := strings.Split(r.Headers["Accept-Encoding"], ",")
	validHeaders := []string{}

	for _, h := range headers {
		h = strings.Trim(h, " ")
		if slices.Contains(supportedCompression, h) {
			validHeaders = append(validHeaders, h)
		}
	}

	if len(validHeaders) == 0 {
		return
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(s.Body)
	if err != nil {
		panic(err)
	}
	if err := zw.Close(); err != nil {
		panic(err)
	}

	compressed := buf.Bytes()
	s.Body = compressed
	s.Headers = []string{
		"Content-Type: text/plain",
		fmt.Sprintf("Content-Length: %d", len(compressed)),
		fmt.Sprintf("Content-Encoding: %s", strings.Join(validHeaders, ",")),
	}
}
