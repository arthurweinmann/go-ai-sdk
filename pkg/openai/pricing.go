package openai

import "fmt"

// order in array: prompt, completion
var PricingPer1000TokensPerModel = map[Model][]float64{
	GPT4_8k:         {0.03, 0.06},
	GPT4_32k:        {0.06, 0.12},
	GPT3_5_turbo_4k: {0.002, 0.002},
	//
	TextDavinci3_4k:    {0.02, 0.02},
	TextDavinci2_4k:    {0.02, 0.02},
	TextDavinci_1_Edit: {0.02, 0.02},
	CodeDavinci2_8k:    {0.02, 0.02},
	//
	Text_Embedding_Ada_2_8k: {0.0004, 0.0004},
}

// Usage Represents the total token usage per request to OpenAI.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (u *Usage) ComputePrice(m Model) float64 {
	if PricingPer1000TokensPerModel[m] == nil {
		panic(fmt.Errorf("model %s is not yet registered in pricing", m))
	}

	return (float64(u.PromptTokens)/1000)*PricingPer1000TokensPerModel[m][0] +
		(float64(u.CompletionTokens)/1000)*PricingPer1000TokensPerModel[m][1]
}
