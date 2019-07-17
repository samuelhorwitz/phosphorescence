package spotifyclient

import (
	"context"
	"io"
)

type cleanupReadCloser struct {
	body   io.ReadCloser
	cancel context.CancelFunc
}

func (rc cleanupReadCloser) Read(p []byte) (int, error) {
	return rc.body.Read(p)
}

func (rc cleanupReadCloser) Close() error {
	rc.cancel()
	return rc.body.Close()
}
