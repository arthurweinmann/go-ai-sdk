package wcohere

import (
	"github.com/cohere-ai/cohere-go"
)

var DefaultClient *cohere.Client

func InitDefaultClient(defaultApiKey string) error {
	var err error
	DefaultClient, err = cohere.CreateClient(defaultApiKey)
	if err != nil {
		return err
	}

	return nil
}
