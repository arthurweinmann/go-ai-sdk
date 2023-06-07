package openai

import (
	"fmt"

	"github.com/arthurweinmann/go-ai-sdk/pkg/openai/pythontool"
)

// GetMaxRemainingTokens uses openai tiktoken to compute the number of tokens and thus you need to have python3 installed
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

	tokencount += 50

	switch m {
	default:
		return 0, fmt.Errorf("model %s not yet supported", m)
	case GPT3_5_turbo_4k:
		return 4096 - tokencount, nil
	case GPT4_8k:
		return 8192 - tokencount, nil
	case GPT4_32k:
		return 32768 - tokencount, nil
	}
}
