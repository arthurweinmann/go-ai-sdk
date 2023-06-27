package uni

type Provider interface{}

type Universal struct {
	promptsPerModel map[string]string
}

type ChatMessage interface {
	GetContent() string
}

type Chat interface {
	AddMessage(ChatMessage) error
}

type Semantic interface{}

type Generative interface{}

type TextAnalyzer interface{}
