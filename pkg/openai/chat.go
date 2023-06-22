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
	Function  MessageRole = "function"
)

type ChatCompletionMessage struct {
	Role MessageRole `json:"role"`

	// The contents of the message. content is required for all messages except assistant messages with function calls.
	Content string `json:"content"`

	// The name of the author of this message. name is required if role is function,
	// and it should be the name of the function whose response is in the content.
	// May contain a-z, A-Z, 0-9, and underscores, with a maximum length of 64 characters.
	Name string `json:"name,omitempty"`

	// Controls how the model responds to function calls.
	// "none" means the model does not call a function, and responds to the end-user.
	// "auto" means the model can pick between an end-user or calling a function.
	// Specifying a particular function via {"name":\ "my_function"} forces the model to call that function.
	// "none" is the default when no functions are present. "auto" is the default if functions are present.
	/* Example:
	FunctionCall: (map[string]interface {}) (len=2) {
	   (string) (len=4) "name": (string) (len=19) "get_current_weather",
	   (string) (len=9) "arguments": (string) (len=60) "{\n  \"format\": \"celsius\",\n  \"location\": \"Glasgow, Scotland\"\n}"
	}
	*/
	FunctionCall any `json:"function_call,omitempty"` // string or map
}

// Under the hood, functions are injected into the system message in a syntax the model has been trained on.
// This means functions count against the model's context limit and are billed as input tokens.
// If running into context limits, we suggest limiting the number of functions or the length of documentation you provide for function parameters.
type ChatCompletionFunction struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or contain underscores and dashes, with a maximum length of 64.
	Name string `json:"name,omitempty"`

	// The description of what the function does.
	Description string `json:"description,omitempty"`

	// The parameters the functions accepts, described as a JSON Schema object.
	Parameters *FunctionParameters `json:"parameters,omitempty"`
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

	// Only if Type == "array"
	Items *FunctionProperty `json:"items,omitempty"`

	// Only if Type == "object"
	Required []string `json:"required,omitempty"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	// Only required if no default api key was initialized
	APIKEY string `json:"-"`

	Model Model `json:"model"`

	Messages  []ChatCompletionMessage  `json:"messages"`
	Functions []ChatCompletionFunction `json:"functions,omitempty"`

	// Set to -1 to let the function automatically compute the maximum number of remaining token in the context
	// window size of the selected model
	// The function returns an error if there are not enough token left for the provided messages and functions
	//
	// Set to -2 to let the function switch between similar models with different maximum context length depending
	// on the token length of the request
	MaxTokens int `json:"max_tokens,omitempty"`

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
		case -2:
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
		}
	}

	err = request("POST", urlSuffix_chatcompletion, req, resp, req.APIKEY)
	if err != nil {
		return nil, err
	}

	resp.Price = resp.Usage.ComputePrice(req.Model)

	return resp, nil
}
