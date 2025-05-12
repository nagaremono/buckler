package internals

import (
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

	s.Headers = append(s.Headers, fmt.Sprintf("Content-Encoding: %s", strings.Join(validHeaders, ",")))
}
