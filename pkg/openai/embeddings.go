package openai

import (
	"fmt"
)

const urlSuffix_embeddings = "v1/embeddings"

type EmbeddingRequest struct {
	// Only required if no default api key was initialized
	APIKEY string `json:"-"`

	Model Model `json:"model"`

	// Input text to get embeddings for, encoded as a string or array of tokens.
	// To get embeddings for multiple inputs in a single request, pass an array of strings or array of token arrays.
	// Each input must not exceed 8192 tokens in length.
	Input interface{} `json:"input"`

	// Only supported in text-embedding-3 and later models
	Dimensions int `json:"dimensions"`

	User string `json:"user,omitempty"`
}

type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
	//

	Price float64 `json:"price,omitempty"`
}

func CreateEmbedding(req *EmbeddingRequest) (*EmbeddingResponse, error) {
	resp := &EmbeddingResponse{}

	if PricingPer1000TokensPerModel[req.Model] == nil {
		return nil, fmt.Errorf("unknown model: %s", req.Model)
	}

	err := request("POST", urlSuffix_embeddings, req, resp, req.APIKEY)
	if err != nil {
		return nil, err
	}

	resp.Price = resp.Usage.ComputePrice(req.Model)

	return resp, nil
}
