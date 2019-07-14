package spotifyclient

import (
	"errors"
	"net"
)

var ErrNoGetBody = errors.New("no GetBody function defined but request has body, cannot clone")
var ErrTimeout = errors.New("resiliency deadline exceeded")

type Error struct {
	err     error
	timeout bool
}

func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) Timeout() bool {
	netErr, ok := e.err.(net.Error)
	return e.timeout || (ok && netErr.Timeout())
}

func (e Error) Temporary() bool {
	netErr, ok := e.err.(net.Error)
	return ok && netErr.Temporary()
}
