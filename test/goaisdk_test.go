package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
	"github.com/arthurweinmann/go-ai-sdk/pkg/wcohere"
)

func TestSetup(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	b, err := os.ReadFile(filepath.Join(wd, "credentials.json"))
	if err != nil {
		panic(fmt.Errorf("%v (you may need to run `go test -v` in the directory of the test file to enable it to find the credentials.json file with your openai api key)", err))
	}

	creds := &struct {
		OpenAIAPIKey string `json:"OPENAI_API_KEY,omitempty"`
		CohereAPIKey string `json:"COHERE_API_KEY,omitempty"`
	}{}

	err = json.Unmarshal(b, creds)
	if err != nil {
		panic(err)
	}

	err = openai.Init(creds.OpenAIAPIKey)
	if err != nil {
		panic(err)
	}

	err = wcohere.InitDefaultClient(creds.CohereAPIKey)
	if err != nil {
		panic(err)
	}
}

func TestCreateChatCompletion(t *testing.T) {
	req := &openai.ChatCompletionRequest{
		Model:       openai.GPT4_128k_Preview,
		Temperature: 0.7,
		Tools: []openai.ChatCompletionToolCall{
			{
				Type: "function",
				Function: &openai.ChatCompletionFunction{
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
				},
			},
		},
		Messages: []openai.ChatCompletionMessage{{
			Role:    "user",
			Content: "what is the weather like today",
		}},
	}

	resp, err := openai.CreateChatCompletion(req)
	if err != nil {
		t.Fatal("Error:", err)
		return
	}

	if resp.Choices[0].Message.Role != "assistant" {
		t.Fatalf("we got an unexpected response message role: %s", resp.Choices[0].Message.Role)
	}

	// Append the response message and a new user message to the request
	req.Messages = append(req.Messages, resp.Choices[0].Message)
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: "I'm in Glasgow, Scotland",
	})

	// Make another request with the updated messages
	resp, err = openai.CreateChatCompletion(req)
	if err != nil {
		t.Fatal("Error:", err)
		return
	}

	if len(resp.Choices[0].Message.ToolCalls) == 0 {
		t.Fatal("unexpected response")
	}
	if resp.Choices[0].Message.ToolCalls[0].Function.Name != "get_current_weather" {
		t.Fatal("unexpected response")
	}
}
