package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/arthurweinmann/go-ai-sdk/internal/requests"
)

var Err429 = errors.New("ratelimit or overload")

const (
	initialDelay  = 30 * time.Second
	maxRetries    = 7
	backoffFactor = 2
)

var baseurl *url.URL
var retryrequester *requests.RequestRetrier

func init() {
	baseurl, _ = url.Parse("https://api.openai.com")

	retryrequester = requests.NewRequestRetrier(initialDelay, maxRetries, backoffFactor)

	retryrequester.Run()
}

func request(method, path string, body, response any, apikey string) error {
	if apikey == "" && defaultAPIKey == "" {
		return fmt.Errorf("we do not have an openai api key defined as default or provided for this request")
	}

	if apikey == "" {
		apikey = defaultAPIKey
	}

	url := baseurl.ResolveReference(&url.URL{Path: path}).String()

	return retryrequester.Request(&requests.Request{
		URL:         url,
		Method:      method,
		Body:        body,
		Response:    response,
		HTTPTimeout: 300 * time.Second,
		Headers: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", apikey)},
		},
		ParseErrBody: func(b []byte, err error, statusCode int, r *requests.Request) error {
			var errRes ErrorResponse
			if len(b) > 0 {
				err = json.Unmarshal(b, &errRes)
			}
			if err != nil || errRes.Error == nil {
				reqErr := RequestError{
					StatusCode: statusCode,
					Err:        fmt.Errorf("%v: %v", err, string(b)), // sometimes when OpenAI nginx fails the error messages is an HTML Page
				}
				err = fmt.Errorf("error requesting %s, %w", url, &reqErr)
			} else {
				errRes.Error.StatusCode = statusCode
				err = fmt.Errorf("error requesting %s, status code: %d, message: %w", url, statusCode, errRes.Error)
			}

			// Openai also has 500 errors sometimes for now so we retry them all until it stabilizes
			// if statusCode == 429 {
			// 	return Err429
			// }

			// We may get context overflow errors until we can figure out a reliable way to compute the number of tokens induced by the functions list and calls
			// features
			isContextLengthOverflow, delta, e := IsErrorContextLengthOverflow(err.Error())
			if e != nil {
				fmt.Println("Could not compute context length overflow correction:", e)
			} else if isContextLengthOverflow {
				fmt.Println("Encountered maximum context length overflow error, decreasing the maxtokens parameters by", delta)
				switch t := r.Body.(type) {
				case *CompletionRequest:
					t.MaxTokens += delta
				case *ChatCompletionRequest:
					t.MaxTokens += delta
				}
			}

			return err
		},
		IsErrorFatal: func(err error) bool {
			// Openai also has 500 errors sometimes for now so we retry them all until it stabilizes
			return false
		},
	})
}
