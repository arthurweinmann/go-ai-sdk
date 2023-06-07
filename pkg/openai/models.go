package openai

type Model string

const (
	Text_Embedding_Ada_2_8k Model = "text-embedding-ada-002"

	/*
		More capable than any GPT-3.5 model, able to do more complex tasks,
		and optimized for chat. Will be updated with our latest model iteration.
	*/
	GPT4_8k Model = "gpt-4"

	/*
		Same capabilities as the base gpt-4 mode but with 4x the context
		length. Will be updated with our latest model iteration.
	*/
	GPT4_32k Model = "gpt-4-32k"

	/*
		Most capable GPT-3.5 model and optimized for chat at 1/10th the cost
		of text-davinci-003. Will be updated with our latest model iteration.
	*/
	GPT3_5_turbo_4k Model = "gpt-3.5-turbo"

	/*
		Can do any language task with better quality, longer output,
		and consistent instruction-following than the curie, babbage, or ada models.  Also supports inserting completions within text.
	*/
	TextDavinci3_4k Model = "text-davinci-003"

	/*
		Similar capabilities to text-davinci-003 but trained with
		supervised fine-tuning instead of reinforcement learning
	*/
	TextDavinci2_4k Model = "text-davinci-002"

	TextDavinci_1_Edit Model = "text-davinci-edit-001"

	/*
		Optimized for code-completion tasks
	*/
	CodeDavinci2_8k Model = "code-davinci-002"
)
