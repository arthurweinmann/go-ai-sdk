package openai

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var reExtractNumbers = regexp.MustCompile(`(?m)(\d+)`)

// APIError provides error information returned by the OpenAI API.
type APIError struct {
	Code       *string `json:"code,omitempty"`
	Message    string  `json:"message"`
	Param      *string `json:"param,omitempty"`
	Type       string  `json:"type"`
	StatusCode int     `json:"-"`
}

// RequestError provides informations about generic request errors.
type RequestError struct {
	StatusCode int
	Err        error
}

type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *RequestError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("status code %d", e.StatusCode)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

// IsErrorContextLengthOverflow returns true and the delta to apply to the MaxTokens request parameter
// if it is a context length overflow, false otherwise
func IsErrorContextLengthOverflow(err string) (bool, int, error) {
	err = strings.ToLower(err)
	if (strings.Contains(err, "maximum") && strings.Contains(err, "context") && strings.Contains(err, "length")) ||
		strings.Contains(err, "reduce") && strings.Contains(err, "length") && strings.Contains(err, "context") {

		var first, second int
		for _, ns := range reExtractNumbers.FindAllString(err, -1) {
			n, err := strconv.Atoi(ns)
			if err != nil {
				panic(fmt.Errorf("Should not happen since we use our own regexp to extract numbers: %v", err))
			}
			if n > first {
				second = first
				first = n
			} else if n > second {
				second = n
			}
		}

		if first == 0 || second == 0 {
			return false, 0, fmt.Errorf("We could not identify the numbers in the context length overflow error")
		}

		return true, second - first - 1, nil
	}

	return false, 0, nil
}
