package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
	"github.com/davecgh/go-spew/spew"
)

var creds = &struct {
	OpenAIAPIKey string `json:"OPENAI_API_KEY,omitempty"`
}{}

func TestFunctionCall(t *testing.T) {
	err := testSetup()
	if err != nil {
		panic(err)
	}
	defer testCleanup()

	err = testFunctionCall()
	if err != nil {
		panic(err)
	}
}

func testSetup() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	b, err := os.ReadFile(filepath.Join(wd, "credentials.json"))
	if err != nil {
		return fmt.Errorf("%v (you may need to run `go test -v` in the directory of the test file to enable it to find the credentials.json file with your openai api key)", err)
	}

	err = json.Unmarshal(b, creds)
	if err != nil {
		return err
	}

	return nil
}

func testCleanup() {
}

func testFunctionCall() error {
	req := &openai.ChatCompletionRequest{
		APIKEY:      creds.OpenAIAPIKey,
		Model:       openai.GPT3_5_turbo_4k_0613,
		Temperature: 0.7,
		Functions: []openai.ChatCompletionFunction{{
			Name:        "get_current_weather",
			Description: "Get the current weather",
			Parameters: &openai.FunctionParameters{
				Type:     "object",
				Required: []string{"location", "format"},
				Properties: map[string]openai.FunctionProperty{
					"location": {
						Type:        "string",
						Description: "The city and state, e.g. San Francisco, CA",
					},
					"format": {
						Type:        "string",
						Description: "The temperature unit to use. Infer this from the users location.",
						Enum:        []string{"celsius", "fahrenheit"},
					},
				},
			},
		}},
		Messages: []openai.ChatCompletionMessage{{
			Role:    "user",
			Content: "what is the weather like today",
		}},
	}
	resp, err := openai.CreateChatCompletion(req)
	if err != nil {
		return err
	}

	spew.Dump(resp.Choices[0])

	req.Messages = append(req.Messages, resp.Choices[0].Message)
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: "I'm in Glasgow, Scotland",
	})
	resp, err = openai.CreateChatCompletion(req)
	if err != nil {
		return err
	}

	spew.Dump(resp.Choices[0])

	return nil
}
