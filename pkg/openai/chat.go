package openai

import (
	"fmt"
)

const urlSuffix_chatcompletion = "v1/chat/completions"

type MessageRole string

const (
	System    MessageRole = "system"
	User      MessageRole = "user"
	Assistant MessageRole = "assistant"
	Tool      MessageRole = "tool"
	// Note: The 'Function' role is deprecated
)

type ContentPart struct {
	Type string `json:"type"`           // "text" or "image_url"
	Text string `json:"text,omitempty"` // for text type
	// For image_url type
	URL    string `json:"url,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type ChatCompletionMessage struct {
	Role      MessageRole `json:"role"`
	Content   interface{} `json:"content"` // string, []ContentPart, or null
	ToolCalls []*ToolCall `json:"tool_calls,omitempty"`
	// Deprecated: Name, FunctionCall
}

type ToolCall struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"` // Currently, only "function" is supported
	Function *Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON format, validation needed
}

type FunctionParameters struct {
	Type       string                      `json:"type"` // object
	Properties map[string]FunctionProperty `json:"properties"`
	Required   []string                    `json:"required,omitempty"`
}

type FunctionProperty struct {
	Type        string   `json:"type"` // string, number, integer, array, boolean, object, null, anyof
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	// For arrays
	Items *FunctionProperty `json:"items,omitempty"`
	// For objects
	Properties map[string]FunctionProperty `json:"properties,omitempty"`
	Required   []string                    `json:"required,omitempty"`
}

type ChatCompletionToolCall struct {
	Type     string                  `json:"type"` // Currently, only "function" is supported
	Function *ChatCompletionFunction `json:"function"`
}

type ChatCompletionFunction struct {
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Parameters  *FunctionParameters `json:"parameters,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type"` // "text" or "json_object"
}

type ChatCompletionRequest struct {
	APIKEY     string `json:"-"`
	MaxRetries int    `json:"-"`

	Model          Model                    `json:"model"`
	Messages       []ChatCompletionMessage  `json:"messages"`
	Tools          []ChatCompletionToolCall `json:"tools,omitempty"`
	ResponseFormat *ResponseFormat          `json:"response_format,omitempty"`
	Seed           *int                     `json:"seed,omitempty"`
	ToolChoice     interface{}              `json:"tool_choice,omitempty"` // string or object
	MaxTokens      int                      `json:"max_tokens,omitempty"`

	Temperature      float32        `json:"temperature,omitempty"`
	TopP             float32        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID                string                 `json:"id"`
	Choices           []ChatCompletionChoice `json:"choices"`
	FinishReason      string                 `json:"finish_reason"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	SystemFingerprint string                 `json:"system_fingerprint"`
	Object            string                 `json:"object"`
	Usage             Usage                  `json:"usage"`

	Price float64 `json:"price,omitempty"`
}

func CreateChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	resp := &ChatCompletionResponse{}

	if PricingPer1000TokensPerModel[req.Model] == nil {
		return nil, fmt.Errorf("unknown model: %s", req.Model)
	}

	var err error

	if req.MaxTokens < 0 {
		switch req.MaxTokens {
		default:
			return nil, fmt.Errorf("We got an invalid MaxTokens parameter of %d, you may use -1 or -2 or a positive value", req.MaxTokens)
		case -1:
			req.MaxTokens, err = GetMaxRemainingTokensChatCompletion(req)
			if err != nil {
				return nil, err
			}

			if req.MaxTokens <= 16 {
				return nil, fmt.Errorf("we do not have enough token left in the context window of the model %s, we have %d token left", req.Model, req.MaxTokens)
			}
		case -2, -3:
			count, err := CountTokensCompletion(req)
			if err != nil {
				return nil, err
			}
			newmodel := req.Model
			maxcontentlength := int(newmodel.GetContextLength())
			isnext := true
			for isnext && maxcontentlength < count+16 {
				isnext, newmodel = newmodel.GetSimilarWithNextContextLength()
				if !isnext {
					return nil, fmt.Errorf("We do not have a model similar to %s with a larger maximum context length", req.Model)
				}
				maxcontentlength = int(newmodel.GetContextLength())
			}
			req.Model = newmodel
			if req.MaxTokens == -2 {
				req.MaxTokens = maxcontentlength - count
			} else if req.MaxTokens == -3 {
				req.MaxTokens = 0
			}
		}
	}

	err = request("POST", urlSuffix_chatcompletion, req, resp, req.APIKEY, req.MaxRetries)
	if err != nil {
		return nil, err
	}

	resp.Price = resp.Usage.ComputePrice(req.Model)

	return resp, nil
}
