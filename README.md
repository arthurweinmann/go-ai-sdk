# Golang AI SDKs

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

    messages := []ChatCompletionMessage{
        {
            Role:    "system",
            Content: "You are a helpful assistant.",
        },
        {
            Role:    "user",
            Content: "Who won the world series in 2022?",
        },
    }

    req := &ChatCompletionRequest{
        APIKEY:   apiKey,
        Model:    model,
        Messages: messages,
        MaxTokens: 60,
    }

    resp, err := CreateChatCompletion(req)
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

# Disclaimer

⚠️ In Active Development ⚠️

Please be aware that this project is currently in the midst of active development.

This project is a fundamental part of a bigger project I am concurrently developing. Therefore, it is very specific to the needs and requirements of that particular project. Some of the modules or features in this project might not be directly usable or understandable outside the context of the bigger project.

Contributors and issue reporting are very welcomed.

Please use this project at your own risk. This project is provided "as is" without any express or implied warranties.

# License

Please refer to the LICENSE file for information on the rights and limitations under the project's license.