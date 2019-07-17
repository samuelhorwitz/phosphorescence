package spotifyclient

import (
	"errors"
	"net"
)

var ErrNoGetBody = errors.New("no GetBody function defined but request has body, cannot clone")

type Error struct {
	err error
}

func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) Timeout() bool {
	netErr, ok := e.err.(net.Error)
	return ok && netErr.Timeout()
}

func (e Error) Temporary() bool {
	netErr, ok := e.err.(net.Error)
	return ok && netErr.Temporary()
}
