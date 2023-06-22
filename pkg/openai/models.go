package openai

type Model string

const (
	Text_Embedding_Ada_2_8k Model = "text-embedding-ada-002"

	GPT4_8k Model = "gpt-4"

	GPT4_8k_0613 Model = "gpt-4-0613"

	GPT4_32k Model = "gpt-4-32k"

	GPT4_32k_0613 Model = "gpt-4-32k-0613"

	GPT3_5_turbo_4k Model = "gpt-3.5-turbo"

	GPT3_5_turbo_16k Model = "gpt-3.5-turbo-16k"

	GPT3_5_turbo_4k_0613 Model = "gpt-3.5-turbo-0613"

	GPT3_5_turbo_16k_0613 Model = "gpt-3.5-turbo-16k-0613"

	GPT3_5_turbo_4k_0301 Model = "gpt-3.5-turbo-0301"

	TextDavinci3_4k Model = "text-davinci-003"

	TextDavinci2_4k Model = "text-davinci-002"

	TextDavinci_1_Edit Model = "text-davinci-edit-001"

	CodeDavinci2_8k Model = "code-davinci-002"
)

type ContextLength int

const (
	Context4K  ContextLength = 4096
	Context8K  ContextLength = 8192
	Context16K ContextLength = 16384
	Context32K ContextLength = 32768
)

func (m Model) GetContextLength() ContextLength {
	switch m {
	default:
		panic("Model does not exist")
	case Text_Embedding_Ada_2_8k:
		return Context8K
	case GPT4_8k, GPT4_8k_0613:
		return Context8K
	case GPT4_32k, GPT4_32k_0613:
		return Context32K
	case GPT3_5_turbo_4k, GPT3_5_turbo_4k_0613, GPT3_5_turbo_4k_0301:
		return Context4K
	case GPT3_5_turbo_16k, GPT3_5_turbo_16k_0613:
		return Context16K
	case TextDavinci3_4k, TextDavinci2_4k, TextDavinci_1_Edit:
		return Context4K
	case CodeDavinci2_8k:
		return Context8K
	}
}

func (m Model) GetSimilarWithNextContextLength() (bool, Model) {
	switch m {
	default:
		panic("Model does not exist")
	case Text_Embedding_Ada_2_8k:
		return false, ""
	case GPT4_8k:
		return true, GPT4_32k
	case GPT4_8k_0613:
		return true, GPT4_32k_0613
	case GPT4_32k, GPT4_32k_0613:
		return false, ""
	case GPT3_5_turbo_4k:
		return true, GPT3_5_turbo_16k
	case GPT3_5_turbo_4k_0301:
		return false, ""
	case GPT3_5_turbo_4k_0613:
		return true, GPT3_5_turbo_16k_0613
	case GPT3_5_turbo_16k, GPT3_5_turbo_16k_0613:
		return false, ""
	case TextDavinci3_4k, TextDavinci2_4k, TextDavinci_1_Edit:
		return false, ""
	case CodeDavinci2_8k:
		return false, ""
	}
}
