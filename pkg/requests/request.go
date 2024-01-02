package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Request struct {
	Method   string
	URL      string
	Body     any
	Response any
	Headers  http.Header

	HTTPTimeout  time.Duration
	ParseErrBody func(body []byte, err error, statusCode int, r *Request) error
}

func Send(r *Request) error {
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

	for hname, hs := range r.Headers {
		for i := 0; i < len(hs); i++ {
			if i == 0 {
				req.Header.Set(hname, hs[i])
			} else {
				req.Header.Add(hname, hs[i])
			}
		}
	}

	req.Header.Set("Accept", "application/json")

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "github.com/arthurweinmann/go-ai-sdk")
	}

	resp, err := (&http.Client{
		Timeout: r.HTTPTimeout,
	}).Do(req)
	if err != nil {
		return fmt.Errorf("http request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
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
