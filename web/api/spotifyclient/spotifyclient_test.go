package spotifyclient_test

import (
	"context"
	"fmt"
	"github.com/samuelhorwitz/phosphorescence/api/spotifyclient"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

type response struct {
	res *http.Response
	err error
}

const okBody = "ok"
const failBody = "fail"

func TestSimple(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(okBody))
	}))
	defer server.Close()
	timeout := 1 * time.Second
	client := &spotifyclient.SpotifyClient{
		Timeout: timeout,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Bad request: %s", err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Could not execute resilient request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Invalid response code %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response body: %s", err)
	}
	if string(body) != okBody {
		t.Fatalf("Incorrect body: %s", body)
	}
}

func TestSimple429(t *testing.T) {
	var responseCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch responseCount {
		case 0:
			w.Header().Set("Retry-After", "1")
			http.Error(w, failBody, http.StatusTooManyRequests)
		default:
			w.Write([]byte(okBody))
		}
		responseCount++
	}))
	defer server.Close()
	timeout := 2 * time.Second
	client := &spotifyclient.SpotifyClient{
		Timeout: timeout,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Bad request: %s", err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Could not execute resilient request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Invalid response code %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response body: %s", err)
	}
	if string(body) != okBody {
		t.Fatalf("Incorrect body: %s", body)
	}
}

func TestPost429RetainsBodyOnSubsequent(t *testing.T) {
	var responseCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Could not read request body: %s", err)
		}
		if string(body) != "FOOBAR" {
			t.Fatalf("Incorrect body: %s", body)
		}
		switch responseCount {
		case 0:
			w.Header().Set("Retry-After", "1")
			http.Error(w, failBody, http.StatusTooManyRequests)
		default:
			w.Write([]byte(okBody))
		}
		responseCount++
	}))
	defer server.Close()
	timeout := 2 * time.Second
	client := &spotifyclient.SpotifyClient{
		Timeout: timeout,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
	req, err := http.NewRequest("POST", server.URL, strings.NewReader("FOOBAR"))
	if err != nil {
		t.Fatalf("Bad request: %s", err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Could not execute resilient request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Invalid response code %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response body: %s", err)
	}
	if string(body) != okBody {
		t.Fatalf("Incorrect body: %s", body)
	}
}

func TestSpotifyClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Millisecond)
		w.Write([]byte(okBody))
	}))
	defer server.Close()
	client := &spotifyclient.SpotifyClient{
		Timeout: 1,
		Client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Bad request: %s", err)
	}
	_, err = client.Do(req)
	t.Log(err)
	if err, ok := err.(net.Error); !ok || !err.Timeout() {
		t.Fatalf("Expected Spotify client timeout hit")
	}
}

func TestHTTPClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Millisecond)
		w.Write([]byte(okBody))
	}))
	defer server.Close()
	client := &spotifyclient.SpotifyClient{
		Timeout: 3 * time.Second,
		Client: &http.Client{
			Timeout: 1,
		},
	}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Bad request: %s", err)
	}
	_, err = client.Do(req)
	if err, ok := err.(net.Error); !ok || !err.Timeout() {
		t.Fatalf("Expected HTTP client timeout hit")
	}
}

func TestExternalCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.Write([]byte(okBody))
	}))
	defer server.Close()
	timeout := 1 * time.Second
	client := &spotifyclient.SpotifyClient{
		Timeout: timeout,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Bad request: %s", err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan error)
	go func() {
		_, err := client.Do(req)
		c <- err
	}()
	cancel()
	err = <-c
	t.Log(err)
	if err, ok := err.(net.Error); !ok || !err.Timeout() {
		t.Fatalf("Expected external cancellation")
	}
}

func TestBarrage(t *testing.T) {
	for i, tester := range []struct {
		limiter       *rate.Limiter
		expectedValid int
		expectedError int
		timeout       time.Duration
	}{
		{
			limiter: rate.NewLimiter(10, 1),
			timeout: 1 * time.Second,
		},
		{
			limiter: rate.NewLimiter(40, 1),
			timeout: 1 * time.Second,
		},
		{
			limiter: rate.NewLimiter(5, 1),
			timeout: 10 * time.Second,
		},
		{
			limiter: rate.NewLimiter(1, 1),
			timeout: 2 * time.Second,
		},
	} {
		fmt.Printf("Loop %d\n", i+1)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reservation := tester.limiter.Reserve()
			delay := reservation.Delay()
			if delay == 0 {
				fmt.Println("[SERVER] Accepted")
				w.Write([]byte(okBody))
			} else {
				ds := uint64(delay.Seconds() + 1)
				fmt.Printf("[SERVER] Retry after %d\n", ds)
				w.Header().Set("Retry-After", fmt.Sprintf("%d", ds))
				http.Error(w, failBody, http.StatusTooManyRequests)
			}
		}))
		defer server.Close()
		client := &spotifyclient.SpotifyClient{
			Timeout: tester.timeout,
			Client: &http.Client{
				Timeout: tester.timeout,
			},
		}
		var wg sync.WaitGroup
		wg.Add(2)
		waitChan := make(chan struct{})
		go func() {
			defer close(waitChan)
			wg.Wait()
		}()
		chan1 := make(chan response)
		chan2 := make(chan response)
		go func() {
			for i := 0; i < 5; i++ {
				fmt.Println("[CLIENT 1] Making request")
				req, err := http.NewRequest("GET", server.URL, nil)
				if err != nil {
					t.Fatalf("Bad request: %s", err)
				}
				res, err := client.Do(req)
				fmt.Println("[CLIENT 1] Request complete")
				chan1 <- response{res, err}
				time.Sleep(250 * time.Millisecond)
			}
			wg.Done()
		}()
		go func() {
			time.Sleep(50 * time.Millisecond)
			for i := 0; i < 5; i++ {
				fmt.Println("[CLIENT 2] Making request")
				req, err := http.NewRequest("GET", server.URL, nil)
				if err != nil {
					t.Fatalf("Bad request: %s", err)
				}
				res, err := client.Do(req)
				fmt.Println("[CLIENT 2] Request complete")
				chan2 <- response{res, err}
				time.Sleep(333 * time.Millisecond)
			}
			wg.Done()
		}()
		validCount := 0
		errorCount := 0
		done := false
		for {
			select {
			case <-waitChan:
				done = true
			case res := <-chan1:
				if res.err != nil {
					errorCount++
					fmt.Printf("[CLIENT 1] Received error: %s\n", res.err)
				} else {
					validCount++
					fmt.Printf("[CLIENT 1] Received response: %v\n", res.res)
				}
			case res := <-chan2:
				if res.err != nil {
					errorCount++
					fmt.Printf("[CLIENT 2] Received error: %s\n", res.err)
				} else {
					validCount++
					fmt.Printf("[CLIENT 2] Received response: %v\n", res.res)
				}
			}
			if done {
				break
			}
		}
		fmt.Printf("Valid %d, Errored %d\n", validCount, errorCount)
	}
}
