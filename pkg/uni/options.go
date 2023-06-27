package uni

import "github.com/arthurweinmann/go-ai-sdk/pkg/openai"

type withOpenAIOption struct {
	APIKey string

	Model openai.Model

	// Set to -1 to let the function automatically compute the maximum number of remaining token in the context
	// window size of the selected model
	// The function returns an error if there are not enough token left for the provided messages and functions
	//
	// Set to -2 to let the function switch between similar models with different maximum context length depending
	// on the token length of the request (for example going automatically from GPT3_5_turbo_4k to GPT3_5_turbo_16k).
	// In this mode, it will set the MaxTokens parameter to the maximum remaining ones.
	//
	// Set to -3 for a similar switching behaviour as -2 but leaves MaxTokens for the default value to apply
	MaxTokens int `json:"max_tokens,omitempty"`

	Temperature float32 `json:"temperature,omitempty"`
	TopP        float32 `json:"top_p,omitempty"`
	N           int     `json:"n,omitempty"`

	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
}

func (*withOpenAIOption) EmbedderOption() {}

type withCohereOption struct {
	APIKey   string
	Model    string
	Truncate string
}

func (*withCohereOption) EmbedderOption() {}

type WithProviderOption interface {
	WithProviderOption()
}

type providerIden32 string

func (providerIden32) WithProviderOption() {}

func (providerIden32) WithFloat32ProviderOption() {}

type providerIden64 string

func (providerIden64) WithProviderOption() {}

func (providerIden64) WithFloat64ProviderOption() {}

func WithOpenAI() providerIden32 {
	return providerIden32("openai")
}

func WithCohere() providerIden64 {
	return providerIden64("cohere")
}
