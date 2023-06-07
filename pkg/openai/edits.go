package openai

import (
	"fmt"
)

const urlSuffix_edits = "v1/edits"

// EditsRequest represents a request structure for Edits API.
type EditsRequest struct {
	// Only required if no default api key was initialized
	APIKEY string `json:"apikey"`

	Model       Model   `json:"model,omitempty"`
	Input       string  `json:"input,omitempty"`
	Instruction string  `json:"instruction,omitempty"`
	N           int     `json:"n,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	TopP        float32 `json:"top_p,omitempty"`
}

// EditsChoice represents one of possible edits.
type EditsChoice struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}

// EditsResponse represents a response structure for Edits API.
type EditsResponse struct {
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Usage   Usage         `json:"usage"`
	Choices []EditsChoice `json:"choices"`

	//

	Price float64 `json:"price,omitempty"`
}

func CreateEdit(req *EditsRequest) (*EditsResponse, error) {
	resp := &EditsResponse{}

	if PricingPer1000TokensPerModel[req.Model] == nil {
		return nil, fmt.Errorf("unknown model: %s", req.Model)
	}

	err := request("POST", urlSuffix_edits, req, resp, req.APIKEY)
	if err != nil {
		return nil, err
	}

	resp.Price = resp.Usage.ComputePrice(req.Model)

	return resp, nil
}
