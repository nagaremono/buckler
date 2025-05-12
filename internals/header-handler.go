package internals

import (
	"fmt"
	"net"
	"slices"
)

var supportedCompression = []string{"gzip"}

type HeaderHandler func(c net.Conn, r *Request, s *Response)

func CompressionHandler(c net.Conn, r *Request, s *Response) {
	compression := r.Headers["Accept-Encoding"]
	if !slices.Contains(supportedCompression, compression) {
		return
	}

	s.Headers = append(s.Headers, fmt.Sprintf("Content-Encoding: %s", compression))
}
