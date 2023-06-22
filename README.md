<div style="display: block; width: 100%; text-align: center;">
    <img src="doc/banner.png" style="display: inline-block; max-width: 100%;"/>
</div>

A comprehensive collection of Golang SDKs for various AI APIs. Use each independently, or leverage the universal interface to seamlessly interact with multiple APIs in a uniform way. Enables easy swapping between AI services without changing your core codebase.

# Current Status

Currently, we only support OpenAI, and the implementation of its APIs remains under development.

# How to use

Import this golang module with `go get -u github.com/arthurweinmann/go-ai-sdk`.

You may initialize OpenAI's sdk with a default API key. It is optional:

```go
package main

import (
    "github.com/arthurweinmann/go-ai-sdk/pkg/openai"
)

func main() {
    err := openai.Init("YOUR_DEFAULT_API_KEY")
    if err != nil {
        panic(err)
    }
}
```

Here is an example on how to use the new functions calls in the chat completion of OpenAI:

```go
package main

import (
	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
	"github.com/davecgh/go-spew/spew"
)

func main() {
    req := &openai.ChatCompletionRequest{
		APIKEY:      "YOUR_API_KEY",
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
		panic(err)
	}

	spew.Dump(resp.Choices[0])
    /*
(openai.ChatCompletionChoice) {
 Index: (int) 0,
 Message: (openai.ChatCompletionMessage) {
  Role: (openai.MessageRole) (len=9) "assistant",
  Content: (string) (len=22) "Where are you located?",
  Name: (string) "",
  FunctionCall: (interface {}) <nil>
 },
 FinishReason: (string) (len=4) "stop"
}
    */

	req.Messages = append(req.Messages, resp.Choices[0].Message)
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: "I'm in Glasgow, Scotland",
	})
	resp, err = openai.CreateChatCompletion(req)
	if err != nil {
		panic(err)
	}

	spew.Dump(resp.Choices[0])
    /*
(openai.ChatCompletionChoice) {
 Index: (int) 0,
 Message: (openai.ChatCompletionMessage) {
  Role: (openai.MessageRole) (len=9) "assistant",
  Content: (string) "",
  Name: (string) "",
  FunctionCall: (map[string]interface {}) (len=2) {
   (string) (len=4) "name": (string) (len=19) "get_current_weather",
   (string) (len=9) "arguments": (string) (len=60) "{\n  \"format\": \"celsius\",\n  \"location\": \"Glasgow, Scotland\"\n}"
  }
 },
 FinishReason: (string) (len=13) "function_call"
}
    */
}
```

Here is an example on how to create a chat completion:

```go
package main

import (
	"fmt"
	"log"
    "github.com/arthurweinmann/go-ai-sdk/pkg/openai"
)

func main() {
    apiKey := "YOUR_DEFAULT_API_KEY"
    model := openai.GPT4_8k // Or the model you want to use

    messages := []openai.ChatCompletionMessage{
        {
            Role:    "system",
            Content: "You are a helpful assistant.",
        },
        {
            Role:    "user",
            Content: "Who won the world series in 2022?",
        },
    }

    req := &openai.ChatCompletionRequest{
        APIKEY:   apiKey,
        Model:    model,
        Messages: messages,
        MaxTokens: 60,
    }

    resp, err := openai.CreateChatCompletion(req)
    if err != nil {
        log.Fatalf("Failed to create chat completion: %v", err)
    }

    fmt.Printf("Response ID: %s\n", resp.ID)
    fmt.Printf("Created At: %d\n", resp.Created)
    fmt.Printf("Model: %s\n", resp.Model)
    for _, choice := range resp.Choices {
        fmt.Printf("Choice Index: %d\n", choice.Index)
        fmt.Printf("Choice Content: %s\n", choice.Message.Content)
        fmt.Printf("Choice Finish Reason: %s\n", choice.FinishReason)
    }
    fmt.Printf("Total Tokens: %d\n", resp.Usage.TotalTokens)
    fmt.Printf("Prompt Tokens: %d\n", resp.Usage.PromptTokens)
    fmt.Printf("Completion Tokens: %d\n", resp.Usage.CompletionTokens)

    // The total price of your request will be computed and added to the response
    fmt.Printf("Price: %f\n", resp.Price)
}
```

To compute the number of tokens remaining in order to set the max tokens parameter, you may use this helper which calls into python package tiktoken under the hood:

```go
package main

import (
	"fmt"
	"log"
    	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
)

func main() {
    remainingTokens, err := openai.GetMaxRemainingTokens("Who won the world series in 2022?", openai.GPT3_5_turbo_4k)
    if err != nil {
        log.Fatalf("Failed to compute maximum remaining tokens: %v", err)
    }
    fmt.Println("Maximum Remaining Tokens:", remainingTokens)
}
```

## Retry feature

If a request fails, it is added to a waiting list. The error is printed, and the function waits for the retry result asynchronously through a golang channel. A goroutine wakes up every so often and check all the requests in the waiting list. It will pick up those requests whose RetryTime is past the current time and retry them one by one.
    
If a retry fails, it will multiply the delay by backoffFactor (which is 2) and set a new RetryTime for the next retry. Then the request is added back to requestswaiting for the next retry. It ceases to retry requests if one fails, and returns to sleep mode until its next scheduled awakening.

The delay keeps increasing until the request succeeds or until it reaches the maximum number of retries.

# License

Please refer to the LICENSE file for information on the rights and limitations under the project's license.
