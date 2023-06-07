package openai

import (
	"fmt"
)

var defaultAPIKey string

func Init(defaultApiKey string) error {
	_, err := ListModels(&ListModelsRequest{APIKEY: defaultApiKey})
	if err != nil {
		return fmt.Errorf("openai test call failed: %s", err)
	}

	defaultAPIKey = defaultApiKey

	return nil
}
