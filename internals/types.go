package internals

import (
	"fmt"
	"strings"
)

type Request struct {
	Method   string
	Target   string
	Protocol string
	Headers  map[string]string
	Body     []byte
}

type Response struct {
	Status     int
	StatusText string
	Protocol   string
	Headers    []string
	Body       []byte
}

func (r *Response) String() string {
	headers := strings.Join(r.Headers, "\r\n")
	body := string(r.Body)
	s := fmt.Sprintf(
		"%s %d %s\r\n%s\r\n\r\n%s\r\n\r\n",
		r.Protocol, r.Status, r.StatusText, headers, body,
	)
	return s
}

func (r *Response) Bytes() []byte {
	return []byte(r.String())
}
