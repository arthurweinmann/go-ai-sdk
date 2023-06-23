package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RequestRetrier struct {
	initialDelay  time.Duration
	maxRetries    int
	backoffFactor int

	requestswaiting   []*Request
	requestswaitingMu *sync.Mutex
}

type Request struct {
	Method   string
	URL      string
	Body     any
	Response any
	Headers  http.Header

	HTTPTimeout  time.Duration
	ParseErrBody func(body []byte, err error, statusCode int, r *Request) error
	IsErrorFatal func(error) bool

	errCh     chan error
	retryTime int64
	newDelay  time.Duration
}

func NewRequestRetrier(initialDelay time.Duration, maxRetries, backoffFactor int) *RequestRetrier {
	if initialDelay <= 0 {
		initialDelay = 30 * time.Second
	}
	if maxRetries <= 0 {
		maxRetries = 7
	}
	if backoffFactor <= 0 {
		backoffFactor = 2
	}

	return &RequestRetrier{
		initialDelay:      initialDelay,
		backoffFactor:     backoffFactor,
		maxRetries:        maxRetries,
		requestswaitingMu: &sync.Mutex{},
	}
}

func (rr *RequestRetrier) Run() {
	go func() {
	Main:
		for {
			time.Sleep(rr.initialDelay)

			var reqtodo []*Request

			now := time.Now().Unix()

			rr.requestswaitingMu.Lock()
			for i := len(rr.requestswaiting) - 1; i > -1; i-- {
				if rr.requestswaiting[i].retryTime < now {
					reqtodo = append(reqtodo, rr.requestswaiting[i])
					rr.requestswaiting[i] = rr.requestswaiting[len(rr.requestswaiting)-1]
					rr.requestswaiting = rr.requestswaiting[:len(rr.requestswaiting)-1]
				}
			}
			rr.requestswaitingMu.Unlock()

			var retryall = func(startingindex int) {
				for i := startingindex; i < len(reqtodo); i++ {
					r := reqtodo[i]

					r.newDelay *= time.Duration(rr.backoffFactor)
					r.retryTime = time.Now().Add(r.newDelay).Unix()
				}

				if startingindex < len(reqtodo) {
					rr.requestswaitingMu.Lock()
					rr.requestswaiting = append(rr.requestswaiting, reqtodo[startingindex:]...)
					rr.requestswaitingMu.Unlock()
				}
			}

			for i := 0; i < len(reqtodo); i++ {
				r := reqtodo[i]
				err := rr.requestnowait(r)
				if err != nil {
					if r.IsErrorFatal(err) {
						r.errCh <- err
						time.Sleep(10*time.Second + time.Duration(rand.Intn(10)))
						continue
					}

					retryall(i)
					continue Main
				}

				r.errCh <- nil

				time.Sleep(time.Duration(float64(rand.Intn(100)) / 100.0 * float64(rr.initialDelay)))
			}
		}
	}()
}

func (rr *RequestRetrier) Request(r *Request) error {
	err := rr.requestnowait(r)
	if err != nil {
		if r.IsErrorFatal(err) {
			return err
		}

		r.errCh = make(chan error, 1)
		r.retryTime = time.Now().Add(rr.initialDelay).Unix()
		r.newDelay = rr.initialDelay * time.Duration(rr.backoffFactor)

		fmt.Printf("* Error: %s, retrying in %v...\n", err, rr.initialDelay)

		rr.requestswaitingMu.Lock()
		rr.requestswaiting = append(rr.requestswaiting, r)
		rr.requestswaitingMu.Unlock()

		return <-r.errCh
	}

	return nil
}

func (rr *RequestRetrier) requestnowait(r *Request) error {
	var err error

	var jsbody []byte
	if r.Body != nil {
		switch r.Body.(type) {
		case nil:
		default:
			jsbody, err = json.Marshal(r.Body)
			if err != nil {
				return err
			}
		}
	}

	var req *http.Request

	if jsbody != nil {
		req, err = http.NewRequest(r.Method, r.URL, bytes.NewReader(jsbody))
	} else {
		req, err = http.NewRequest(r.Method, r.URL, nil)
	}
	if err != nil {
		return fmt.Errorf("http request: %v", err)
	}

	if len(jsbody) > 0 {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Length", strconv.Itoa(len(jsbody)))
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "github.com/arthurweinmann/go-ai-sdk")

	for hname, hs := range r.Headers {
		for i := 0; i < len(hs); i++ {
			if i == 0 {
				req.Header.Set(hname, hs[i])
			} else {
				req.Header.Add(hname, hs[i])
			}
		}
	}

	resp, err := (&http.Client{
		Timeout: r.HTTPTimeout,
	}).Do(req)
	if err != nil {
		return fmt.Errorf("http request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		b, errReadBody := io.ReadAll(resp.Body)
		if errReadBody != nil {
			err = fmt.Errorf("We received %v and could not read the request's body: %v", err, errReadBody)
		}

		err = r.ParseErrBody(b, err, resp.StatusCode, r)

		return err
	}

	err = json.NewDecoder(resp.Body).Decode(r.Response)
	if err != nil {
		return fmt.Errorf("unmarshal response: %v", err)
	}

	return nil
}
