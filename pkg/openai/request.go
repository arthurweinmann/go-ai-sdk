package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var Err429 = errors.New("ratelimit or overload")

const (
	initialDelay  = 30 * time.Second
	maxRetries    = 7
	backoffFactor = 2
)

type requestIn429 struct {
	method   string
	path     string
	body     any
	response any
	apikey   string

	ch chan error

	RetryTime int64
	NewDelay  time.Duration
}

var requestswaiting = []*requestIn429{}
var requestswaitingMu sync.Mutex

var baseurl *url.URL
var httpclient *http.Client

func init() {
	baseurl, _ = url.Parse("https://api.openai.com")

	httpclient = &http.Client{
		Timeout: 300 * time.Second,
	}

	go func() {
	Main:
		for {
			time.Sleep(initialDelay)

			var reqtodo []*requestIn429

			now := time.Now().Unix()

			requestswaitingMu.Lock()
			for i := len(requestswaiting) - 1; i > -1; i-- {
				if requestswaiting[i].RetryTime < now {
					reqtodo = append(reqtodo, requestswaiting[i])
					requestswaiting[i] = requestswaiting[len(requestswaiting)-1]
					requestswaiting = requestswaiting[:len(requestswaiting)-1]
				}
			}
			requestswaitingMu.Unlock()

			var retryall = func(startingindex int) {
				for i := startingindex; i < len(reqtodo); i++ {
					r := reqtodo[i]

					r.NewDelay *= backoffFactor
					r.RetryTime = time.Now().Add(r.NewDelay).Unix()
				}

				if startingindex < len(reqtodo) {
					requestswaitingMu.Lock()
					requestswaiting = append(requestswaiting, reqtodo[startingindex:]...)
					requestswaitingMu.Unlock()
				}
			}

			for i := 0; i < len(reqtodo); i++ {
				r := reqtodo[i]
				err := requestnowait(r.method, r.path, r.body, r.response, r.apikey)
				if err != nil {
					// Openai also has 500 errors sometimes for now
					// if err != Err429 {
					// 	r.ch <- err
					// 	time.Sleep(10*time.Second + time.Duration(rand.Intn(10)))
					// 	continue
					// }

					retryall(i)
					continue Main
				}

				r.ch <- nil

				time.Sleep(time.Duration(float64(rand.Intn(100)) / 100.0 * float64(initialDelay)))
			}
		}
	}()
}

func request(method, path string, body, response any, apikey string) error {
	if apikey == "" && defaultAPIKey == "" {
		return fmt.Errorf("we do not have an openai api key defined as default or provided for this request")
	}

	err := requestnowait(method, path, body, response, apikey)
	if err != nil {
		// Openai also has 500 errors sometimes for now
		// if err != Err429 {
		// 	return err
		// }

		r := &requestIn429{
			method:    method,
			path:      path,
			body:      body,
			response:  response,
			apikey:    apikey,
			ch:        make(chan error, 1),
			RetryTime: time.Now().Add(initialDelay).Unix(),
			NewDelay:  initialDelay * backoffFactor,
		}

		fmt.Printf("* Error: %s, retrying in %v...\n", err, initialDelay)

		requestswaitingMu.Lock()
		requestswaiting = append(requestswaiting, r)
		requestswaitingMu.Unlock()

		return <-r.ch
	}

	return nil
}

func requestnowait(method, path string, body, response any, apikey string) error {
	if apikey == "" {
		apikey = defaultAPIKey
	}

	var err error

	rel := &url.URL{Path: path}
	u := baseurl.ResolveReference(rel)
	var jsbody []byte
	if body != nil {
		jsbody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	var req *http.Request

	if jsbody != nil {
		req, err = http.NewRequest(method, u.String(), bytes.NewReader(jsbody))
	} else {
		req, err = http.NewRequest(method, u.String(), nil)
	}
	if err != nil {
		return fmt.Errorf("http request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "github.com/arthurweinmann/go-ai-sdk")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))

	resp, err := httpclient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		b, err := io.ReadAll(resp.Body)
		if err == nil {
			err = json.Unmarshal(b, &errRes)
		}
		if err != nil || errRes.Error == nil {
			reqErr := RequestError{
				StatusCode: resp.StatusCode,
				Err:        fmt.Errorf("%v: %v", err, string(b)), // sometimes when OpenAI nginx fails the error messages is an HTML Page
			}
			err = fmt.Errorf("error requesting %s, %w", u, &reqErr)
		} else {
			errRes.Error.StatusCode = resp.StatusCode
			err = fmt.Errorf("error requesting %s, status code: %d, message: %w", u, resp.StatusCode, errRes.Error)
		}

		// if resp.StatusCode == 429 {
		// 	utils.Debug("\n    * Error: %s\n", err)
		// 	return Err429
		// }

		// We may get context overflow errors until we can figure out a reliable way to compute the number of tokens induced by the functions list and calls
		// features
		isContextLengthOverflow, delta, e := IsErrorContextLengthOverflow(err.Error())
		if e != nil {
			fmt.Println("Could not compute context length overflow correction:", e)
		} else if isContextLengthOverflow {
			fmt.Println("Encountered maximum context length overflow error, decreasing the maxtokens parameters by", delta)
			switch t := body.(type) {
			default:
				panic(fmt.Errorf("Should not happen: %T", t))
			case *CompletionRequest:
				t.MaxTokens += delta
			case *ChatCompletionRequest:
				t.MaxTokens += delta
			case *EditsRequest:
			case *EmbeddingRequest:
			}
		}

		return err
	}

	// b, _ := io.ReadAll(resp.Body)
	// spew.Dump(string(b))
	// err = json.Unmarshal(b, response)

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return fmt.Errorf("unmarshal response: %v", err)
	}

	return nil
}
