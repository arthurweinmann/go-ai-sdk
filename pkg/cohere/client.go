package cohere

import "github.com/cohere-ai/cohere-go"

func NewClient(apikey string) (*cohere.Client, error) {
	return cohere.CreateClient(apikey)
}
