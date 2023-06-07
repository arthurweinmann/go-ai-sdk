package openai

import (
	"fmt"
)

const urlSuffix_completion = "v1/completions"

// CompletionRequest represents a request structure for completion API.
type CompletionRequest struct {
	// Only required if no default api key was initialized
	APIKEY string `json:"apikey"`

	Model Model `json:"model"`

	// The prompt(s) to generate completions for, encoded as a string, array of strings, array of tokens, or array of token arrays.
	Prompt any `json:"prompt,omitempty"`

	Suffix           string         `json:"suffix,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      float32        `json:"temperature,omitempty"`
	TopP             float32        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	LogProbs         int            `json:"logprobs,omitempty"`
	Echo             bool           `json:"echo,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	BestOf           int            `json:"best_of,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

// CompletionChoice represents one of possible completions.
type CompletionChoice struct {
	Text         string        `json:"text"`
	Index        int           `json:"index"`
	FinishReason string        `json:"finish_reason"`
	LogProbs     LogprobResult `json:"logprobs"`
}

// LogprobResult represents logprob result of Choice.
type LogprobResult struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float32            `json:"token_logprobs"`
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

// CompletionResponse represents a response structure for completion API.
type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`

	//

	Price float64 `json:"price,omitempty"`
}

func CreateCompletion(req *CompletionRequest) (*CompletionResponse, error) {
	resp := &CompletionResponse{}

	if PricingPer1000TokensPerModel[req.Model] == nil {
		return nil, fmt.Errorf("unknown model: %s", req.Model)
	}

	err := request("POST", urlSuffix_completion, req, resp, req.APIKEY)
	if err != nil {
		return nil, err
	}

	resp.Price = resp.Usage.ComputePrice(req.Model)

	return resp, nil
}
