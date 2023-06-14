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
