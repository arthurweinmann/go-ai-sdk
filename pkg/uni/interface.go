package uni

type Model interface{}

type Universal interface {
	// SetPrompt enables you to use different prompts for different models
	SetPrompt(Model, string)
}

type ChatMessage interface{}

type Chat interface {
	AddMessage(ChatMessage) error
}

type Semantic interface{}

type Generative interface{}

type TextAnalyzer interface{}
