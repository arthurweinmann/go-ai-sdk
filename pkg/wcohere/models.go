package wcohere

// Command models
const (
	// CommandLight is a smaller, faster version of Command.
	// It's almost as capable but a lot faster.
	// Max Tokens: 4096. Endpoint: Co.generate()
	CommandLight string = "command-light"

	// Command is an instruction-following conversational model that performs
	// language tasks with high quality, more reliably and with a longer context
	// than base generative models.
	// Max Tokens: 4096. Endpoint: Co.generate()
	Command string = "command"
)

// Generation models
const (
	// BaseLight is a smaller, faster version of Base.
	// It's almost as capable but a lot faster.
	// Max Tokens: 2048. Endpoint: Co.generate()
	BaseLight string = "base-light"

	// Base is a model that performs generative language tasks.
	// Max Tokens: 2048. Endpoint: Co.generate()
	Base string = "base"
)

// Representation models
const (
	// EmbedEnglishLightV2 is a smaller, faster version of EmbedEnglishV2.
	// It's almost as capable but a lot faster. It supports English only.
	// Max Tokens: 512. Similarity Metric: Cosine Similarity.
	// Endpoints: Co.Classify(), Co.Embed(), Co.Detect_language(), Co.Tokenize(), Co.Detokenize()
	EmbedEnglishLightV2 string = "embed-english-light-v2.0"

	// EmbedEnglishV2 is a model that allows for text to be classified or turned into embeddings.
	// It supports English only.
	// Max Tokens: 512. Similarity Metric: Cosine Similarity.
	// Endpoints: Co.Classify(), Co.Embed(), Co.Detect_language(), Co.Tokenize(), Co.Detokenize()
	EmbedEnglishV2 string = "embed-english-v2.0"

	// EmbedMultilingualV2 provides multilingual classification and embedding support.
	// Max Tokens: 256. Similarity Metric: Dot Product Similarity.
	// Endpoints: Co.Classify(), Co.Embed(), Co.Detect_language(), Co.Tokenize(), Co.Detokenize()
	EmbedMultilingualV2 string = "embed-multilingual-v2.0"
)

// Rerank models
const (
	// RerankEnglishV2 is a model that allows for re-ranking English language documents.
	// No token limit since it accepts full strings rather than tokens.
	// Endpoint: Co.rerank()
	RerankEnglishV2 string = "rerank-english-v2.0"

	// RerankMultilingualV2 is a model for documents that are not in English.
	// Supports the same languages as EmbedMultilingualV2.
	// No token limit since it accepts full strings rather than tokens.
	// Endpoint: Co.rerank()
	RerankMultilingualV2 string = "rerank-multilingual-v2.0"
)

// Summarize models
const (
	// SummarizeMedium is a smaller, faster version of SummarizeXLarge.
	// It's almost as capable but a lot faster.
	// Max Tokens: 2048. Endpoint: Co.summarize()
	SummarizeMedium string = "summarize-medium"

	// SummarizeXLarge is a model that takes a piece of text and generates a summary.
	// Max Tokens: 2048. Endpoint: Co.summarize()
	SummarizeXLarge string = "summarize-xlarge"
)
