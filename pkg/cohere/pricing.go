package cohere

import "fmt"

func GetEmbedRequestPrice(numTokens int) float64 {
	return float64(numTokens) * 0.0000004
}

func GetGenerateRequestPrice(numTokens int, model string) (float64, error) {
	switch model {
	case "default":
		return float64(numTokens) * 0.000015, nil
	case "custom":
		return float64(numTokens) * 0.000030, nil
	default:
		return 0, fmt.Errorf("We support either `default` or `custom` as a model, but we got: %s", model)
	}

}

func GetSummarizeRequestPrice(numTokens int) float64 {
	return float64(numTokens) * 0.000015
}

// Cohere counts a single search unit as a query with up to 100 documents to be ranked.
// Documents longer than 510 tokens when including the length of the search query will be split up into multiple chunks,
// where each chunk counts as a singular document.
func GetRerankPrice(numSearchUnit int) float64 {
	return float64(numSearchUnit) * 0.001
}
