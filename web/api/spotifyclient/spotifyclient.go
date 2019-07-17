package spotifyclient

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

const defaultRetryAfterSeconds = 5

type SpotifyClient struct {
	Timeout time.Duration
	Client  *http.Client
	wg      waitGroupCond
}

func (c *SpotifyClient) Do(baseReq *http.Request) (res *http.Response, err error) {
	// Set a timeout around the multiple request resiliency.
	baseCtx, cancel := context.WithTimeout(baseReq.Context(), c.Timeout)
	// We cannot just `defer cancel()` as normal because the Body
	// being read is part of the context deadline in an http.Response.
	// Instead, we need to have a slightly more complicated defer
	// which immediately cancels on no body or no response, and
	// otherwise overwrites the io.ReadCloser with our own which
	// cancels on Close (something all good callers should be doing
	// already for http requests).
	defer func() {
		if res == nil || res.Body == nil {
			cancel()
		} else {
			res.Body = cleanupReadCloser{res.Body, cancel}
		}
	}()
	// Build our response and error channels.
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)
	// Create a done signal channel for alerting any Go routines that
	// this function has completed and they should clean up.
	done := make(chan struct{})
	defer close(done)
	// Send all requests in a Go routine and respond in channels.
	go func() {
		for {
			// We want don't want this routine to dangle forever.
			// The waitgroup will be around for other requests,
			// this request might timeout and we won't care anymore
			// about waiting. Therefore we use a channel to handle
			// the wait result so that it won't block if another
			// finishing event is selected on.
			waitChan := make(chan struct{})
			go func() {
				defer close(waitChan)
				c.wg.Wait()
			}()
			// This select is only to block continuing until we know
			// that we can resend or this resilient request was
			// cancelled. We don't want to do anything but block
			// until one of these conditions are met. If they both
			// happen at the same time, Go will randomly choose one,
			// so we don't want to handle cancellation here since it
			// might be discarded due to bad luck.
			select {
			case <-waitChan:
			case <-done:
			}
			// Here we actually handle cancellation. We got past the
			// block above, now let's see if it was cancellation that
			// triggered our continuation. If so, return and end this
			// Go routine. Otherwise, keep going.
			select {
			case <-done:
				return
			default:
			}
			// Clone the original request.
			req := baseReq.Clone(baseCtx)
			// Set the body if one exists. If there's a body but not
			// a GetBody function, we error. A lot of builtins will
			// automatically set GetBody but the user will have to
			// set it for custom readers.
			if req.GetBody != nil {
				var err error
				req.Body, err = baseReq.GetBody()
				if err != nil {
					errorChan <- err
					return
				}
			} else if req.Body != nil {
				errorChan <- ErrNoGetBody
				return
			}
			// Send off the request.
			res, err := c.Client.Do(req)
			// Did we get cancelled in the time it took for the request
			// to resolve or error? If so, let's exit early and end this
			// Go routine.
			select {
			case <-done:
				// We were cancelled, so let's do some safety checks in
				// case error (which we don't care about now) was set
				// and res or it's Body was nil (Body may be nil if error,
				// unsure if that could happen, but whatever).
				if res != nil && res.Body != nil {
					// We can close the response body so we should, we don't
					// want it anymore.
					res.Body.Close()
				}
				return
			default:
			}
			// We got an error, so let's send it off and return. We are not
			// resilient when it comes to Go-level errors, only to 429 HTTP
			// back-off requests.
			if err != nil {
				errorChan <- err
				return
			}
			// If we've been sending too many requests we close the response body,
			// check out many seconds we should back-off for, and then use our wait
			// group to ensure this client backs off whatever thread or function is
			// calling it. Then, loop. Otherwise, everything was good and we got a
			// response. We send that back and exit the Go routine.
			if res.StatusCode == http.StatusTooManyRequests {
				res.Body.Close()
				retryAfterSeconds, err := strconv.Atoi(res.Header.Get("Retry-After"))
				if err != nil {
					retryAfterSeconds = defaultRetryAfterSeconds
				}
				c.backoffFor(time.Duration(retryAfterSeconds) * time.Second)
			} else {
				responseChan <- res
				return
			}
		}
	}()
	// Wait for a response, an error response, an externally cancelled
	// context, or a resiliency retry timeout.
	select {
	case res = <-responseChan:
		return res, nil
	case err = <-errorChan:
		return nil, Error{err}
	case <-baseCtx.Done():
		return nil, Error{baseCtx.Err()}
	}
}

func (c *SpotifyClient) backoffFor(t time.Duration) {
	c.wg.Increment()
	go c.continueRequestsAfter(t)
}

func (c *SpotifyClient) continueRequestsAfter(t time.Duration) {
	time.Sleep(t)
	c.wg.Done()
}
