<p align="center">
    <img src="doc/banner.png"/>
</p>

A comprehensive collection of Golang SDKs for various AI and Knowledge APIs. Use each independently, or leverage the universal interface to seamlessly interact with multiple APIs in a uniform way. Enables easy swapping between AI services without changing your core codebase.

# Current Status

Currently, we support:

    - OpenAI, and the implementation of its APIs remains under development.

    - Google Cloud Natural Language APIs, also under development.

    - Wikipedia (Wikimedia) Knowledge API

# How to use

Import this golang module with `go get -u github.com/arthurweinmann/go-ai-sdk`.

## OpenAI

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
		
		// Set to any positive value to limit the completion number of tokens in the response
		//
		// Set to -1 to let the function automatically compute the maximum number of remaining token in the context
		// window size of the selected model
		// The function returns an error if there are not enough token left for the provided messages and functions
		//
		// Set to -2 to let the function switch between similar models with different maximum context length depending
		// on the token length of the request (for example going automatically from GPT3_5_turbo_4k to GPT3_5_turbo_16k)
		// In this mode, it will set the MaxTokens parameter to the maximum remaining ones.
		//
		// Set to -3 for a similar switching behaviour as -2 but leaves MaxTokens for the default value to apply
		MaxTokens: -2,

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

The `MaxTokens` value controls the length of the response, i.e., the number of tokens it contains. 

If you set it to any positive number, you can cap the response's size accordingly.

If `MaxTokens` is set to -1, the function will automatically adjust the maximum token count to fit within the context window of the chosen model. However, this might cause an error if the available tokens are insufficient for the messages and operations provided.

By setting `MaxTokens` to -2, the function gains the flexibility to alternate between similar models with varying context lengths, depending on the request's token length. For instance, it can seamlessly transition from GPT3_5_turbo_4k to GPT3_5_turbo_16k, depending on the length of the request. In this mode, it will set the MaxTokens parameter to the maximum remaining ones.

Finally, by setting `MaxTokens` to -3, you get a similar switching behaviour as in -2 mode but it leaves MaxTokens for the default value to apply.

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

## Google Natural Language API

You first have to initialize Google Natural Language's sdk with your API key:

```go
package main

import (
    "github.com/arthurweinmann/go-ai-sdk/pkg/googlenl"
)

func main() {
    err := googlenl.Init("YOUR_API_KEY")
    if err != nil {
        panic(err)
    }
}
```

You can for example extract entities from a text:

```go
package main

import (
    "fmt"
    "log"
    "context"
    "github.com/arthurweinmann/go-ai-sdk/pkg/googlenl"
)

func main() {
    err := googlenl.Init("YOUR_API_KEY")
    if err != nil {
        panic(err)
    }

    text := "Google, headquartered in Mountain View, unveiled the new Android phone at the Consumer Electronic Show. Sundar Pichai is the CEO of Google."

	resp, err := googlenl.AnalyzeEntities(context.Background(), text)
	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}

	// Print the results
	fmt.Printf("Language of the text: %s\n", resp.Language)
	for _, entity := range resp.Entities {
		fmt.Printf("Entity: %s\n", entity.Name)
		fmt.Printf("Type: %v\n", entity.Type)
		fmt.Printf("Salience: %.2f\n", entity.Salience)
		for key, value := range entity.Metadata {
			fmt.Printf("Metadata: %s: %s\n", key, value)
		}
		for _, mention := range entity.Mentions {
			fmt.Printf("Mention: Type %v, Text %s\n", mention.Type, mention.Text.Content)
		}
	}
}
```

## Wikipedia (Wikimedia)

Here is a simple example of how you might use this sdk to query Wikipedia for a specific topic and get the related information. In this case, we are interested in "Artificial Intelligence".

First, make sure to initialize your client:

```go
package main

import (
	"fmt"
	"log"
    "github.com/arthurweinmann/go-ai-sdk/pkg/wikipedia"
)

func main() {
	// Initialize the Wikipedia client
	err := wikipedia.Init()
	if err != nil {
		log.Fatal(err)
	}
	// ...

```

Then, you can retrieve a list of pages that start with your query, in this case "Artificial Intelligence", and limit the results to the first 10:

```go
	// Get the first 10 pages that start with "Artificial Intelligence"
	pages, err := wikipedia.Client.GetPrefixResults("Artificial Intelligence", 10)
	if err != nil {
		log.Fatal(err)
	}
	for _, page := range pages {
		fmt.Printf("Page ID: %d, Page Title: %s, Page URL: %s\n", page.ID, page.Title, page.URL)
	}

```

You can also retrieve the extracts for a list of page titles. Let's get the extracts for the first page from the previous result:

```go
	// Get the extracts for the first page
	extracts, err := wikipedia.Client.GetExtracts([]string{pages[0].Title})
	if err != nil {
		log.Fatal(err)
	}
	for _, extract := range extracts {
		fmt.Printf("Page Title: %s, Extract: %s\n", extract.Meta.Title, extract.Extract)
	}

```

Finally, you can retrieve the categories and sections associated with a specific page ID. Let's do it for the first page from the previous result:

```go
	// Get the categories for the first page
	categories, err := wikipedia.Client.GetCategories(pages[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Page Title: %s, Categories: %v\n", categories.Meta.Title, categories.Categories)

	// Get the sections for the first page
	sections, err := wikipedia.Client.GetSections(pages[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Page Title: %s, Sections: %v\n", sections.Meta.Title, sections.Sections)
}
```

# OpenAI Retry feature

If a request fails, it is added to a waiting list. The error is printed, and the function waits for the retry result asynchronously through a golang channel. A goroutine wakes up every so often and check all the requests in the waiting list. It will pick up those requests whose RetryTime is past the current time and retry them one by one.
    
If a retry fails, it will multiply the delay by backoffFactor (which is 2) and set a new RetryTime for the next retry. Then the request is added back to requestswaiting for the next retry. It ceases to retry requests if one fails, and returns to sleep mode until its next scheduled awakening.

The delay keeps increasing until the request succeeds or until it reaches the maximum number of retries.

## Note

The OpenAI API can be unpredictable. At times, it throws 500 error messages even for valid requests. Therefore, we retry any error-producing request. This is beyond the usual practice of retrying just the 429 rate limit errors. 

Also, there's insufficient documentation at the moment on how to count tokens for the new function calls feature. Due to this, we handle context length overflow errors differently. We parse those errors and automatically adjust the MaxTokens parameter for following attempts. This ensures that operations run smoothly.

# License

Please refer to the LICENSE file for information on the rights and limitations under the project's license.
