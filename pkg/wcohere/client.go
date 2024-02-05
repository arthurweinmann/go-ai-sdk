package wcohere

import (
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
)

func NewClient(apikey string) (*cohereclient.Client, error) {
	return cohereclient.NewClient(cohereclient.WithToken(apikey)), nil
}
