package wcohere

import (
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
)

var DefaultClient *cohereclient.Client

func InitDefaultClient(defaultApiKey string) error {
	DefaultClient = cohereclient.NewClient(cohereclient.WithToken(defaultApiKey))

	return nil
}
