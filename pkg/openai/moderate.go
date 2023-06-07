package openai

import (
	"fmt"
)

var urlSuffix_moderate = "v1/moderations"

type ModerateRequest struct {
	// Only required if no default api key was initialized
	APIKEY string `json:"apikey"`

	Input string `json:"input"`
}

type ModerateResponse struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Results []Result `json:"results"`
}

type Result struct {
	Categories     map[string]bool    `json:"categories"`
	CategoryScores map[string]float64 `json:"category_scores"`
	Flagged        bool               `json:"flagged"`
}

func Moderate(req *ModerateRequest) (*ModerateResponse, error) {
	defer fmt.Println("")

	resp := &ModerateResponse{}

	err := request("POST", urlSuffix_moderate, req, resp, req.APIKEY)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
