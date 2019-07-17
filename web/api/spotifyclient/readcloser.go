package spotifyclient

import (
	"context"
	"io"
)

type readCloser struct {
	body   io.ReadCloser
	cancel context.CancelFunc
}

func (rc readCloser) Read(p []byte) (int, error) {
	return rc.body.Read(p)
}

func (rc readCloser) Close() error {
	rc.cancel()
	return rc.body.Close()
}
