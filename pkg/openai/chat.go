package openai

import (
	"fmt"
)

const urlSuffix_chatcompletion = "v1/chat/completions"

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`

	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name,omitempty"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	// Only required if no default api key was initialized
	APIKEY string `json:"-"`

	Model            Model                   `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	Temperature      float32                 `json:"temperature,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"`
	N                int                     `json:"n,omitempty"`
	Stream           bool                    `json:"stream,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int          `json:"logit_bias,omitempty"`
	User             string                  `json:"user,omitempty"`
}

type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`

	Price float64 `json:"price,omitempty"`
}

func CreateChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	resp := &ChatCompletionResponse{}

	if PricingPer1000TokensPerModel[req.Model] == nil {
		return nil, fmt.Errorf("unknown model: %s", req.Model)
	}

	err := request("POST", urlSuffix_chatcompletion, req, resp, req.APIKEY)
	if err != nil {
		return nil, err
	}

	resp.Price = resp.Usage.ComputePrice(req.Model)

	return resp, nil
}
