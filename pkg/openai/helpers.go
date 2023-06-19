package openai

import (
	"fmt"

	"github.com/arthurweinmann/go-ai-sdk/pkg/openai/pythontool"
)

// GetMaxRemainingTokens uses openai tiktoken to compute the number of tokens and thus you need to have python3 installed
// Watchout for functions definitions which count toward the model context length. I did not find information as to the function syntax
// used by openai to compute its number of tokens. You can probably approximate it.
func GetMaxRemainingTokens(prompt string, m Model) (int, error) {
	var encoding string
	switch m {
	default:
		return 0, fmt.Errorf("model %s not implemented yet for GetMaxRemainingTokens", m)
	case GPT3_5_turbo_4k, GPT4_8k, GPT4_32k:
		encoding = "cl100k_base"
	}

	tokencount, err := pythontool.CountTokens(encoding, prompt)
	if err != nil {
		return 0, err
	}

	// TODO: Consider the following: `every reply is primed with <|start|>assistant<|message|>`
	tokencount += 50

	switch m {
	default:
		return 0, fmt.Errorf("model %s not yet supported", m)
	case GPT3_5_turbo_4k, GPT3_5_turbo_4k_0301, GPT3_5_turbo_4k_0613:
		return int(Context4K) - tokencount, nil
	case GPT3_5_turbo_16k_0613, GPT3_5_turbo_16k:
		return int(Context16K) - tokencount, nil
	case GPT4_8k, GPT4_8k_0613:
		return int(Context8K) - tokencount, nil
	case GPT4_32k, GPT4_32k_0613:
		return int(Context32K) - tokencount, nil
	}
}
