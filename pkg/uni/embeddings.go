package uni

import (
	"fmt"
	"sync"

	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
	"github.com/arthurweinmann/go-ai-sdk/pkg/wcohere"
	"github.com/cohere-ai/cohere-go"
)

type Embedder struct {
	err       error
	providers []EmbedderOption
}

type Embedding struct {
	byprovider32 map[string][]float32
	byprovider64 map[string][]float64
}

type EmbedderOption interface {
	EmbedderOption()
}

func WithOpenAIEmbed(model openai.Model, apikeyOptional string) *withOpenAIOption {
	return &withOpenAIOption{
		APIKey: apikeyOptional,
		Model:  model,
	}
}

func WithCohereEmbed(model, truncate, apikeyOptional string) *withCohereOption {
	return &withCohereOption{
		APIKey:   apikeyOptional,
		Model:    model,
		Truncate: truncate,
	}
}

func NewEmbedder(opts ...EmbedderOption) *Embedder {
	emb := &Embedder{}

	if len(opts) == 0 {
		emb.err = fmt.Errorf("We need at least one provider of embeddings")
		return emb
	}

	for i := 0; i < len(opts); i++ {
		switch t := opts[i].(type) {
		default:
			panic(fmt.Errorf("Should not happen: %T", t))
		case *withOpenAIOption:
			emb.providers = append(emb.providers, t)
		case *withCohereOption:
			emb.providers = append(emb.providers, t)
		}
	}

	return emb
}

func (m *Embedder) BatchEmbed(texts []string, opts ...WithProviderOption) ([]*Embedding, error) {
	if m.err != nil {
		return nil, m.err
	}

	useOpenAI := true
	useCohere := true

	if len(opts) > 0 {
		useOpenAI = false
		useCohere = false

		for i := 0; i < len(opts); i++ {
			switch t := opts[i].(type) {
			default:
				panic(fmt.Errorf("Should not happen: %T", t))
			case providerIden:
				switch t {
				default:
					panic(fmt.Errorf("Should not happen: %s", t))
				case "openai":
					useOpenAI = true
				case "cohere":
					useCohere = true
				}
			}
		}
	}

	ret := make([]*Embedding, len(texts))
	for i := 0; i < len(ret); i++ {
		ret[i] = &Embedding{
			byprovider32: map[string][]float32{},
			byprovider64: map[string][]float64{},
		}
	}

	var wg sync.WaitGroup
	var errs []error
	var mu sync.Mutex

	for _, prov := range m.providers {
		switch t := prov.(type) {
		default:
			panic(fmt.Errorf("Should not happen: %T", t))
		case *withOpenAIOption:
			if useOpenAI {
				wg.Add(1)
				go func() {
					defer wg.Done()
					resp, err := openai.CreateEmbedding(&openai.EmbeddingRequest{
						APIKEY: t.APIKey,
						Model:  t.Model,
						Input:  texts,
					})
					mu.Lock()
					defer mu.Unlock()
					if err != nil {
						errs = append(errs, fmt.Errorf("OpenAI: %v", err))
						return
					}
					for i := 0; i < len(resp.Data); i++ {
						ret[resp.Data[i].Index].byprovider32["openai"] = resp.Data[i].Embedding
					}
				}()
			}
		case *withCohereOption:
			if useCohere {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var client *cohere.Client
					if t.APIKey != "" {
						var err error
						client, err = cohere.CreateClient(t.APIKey)
						if err != nil {
							mu.Lock()
							defer mu.Unlock()
							errs = append(errs, fmt.Errorf("Cohere: %v", err))
							return
						}
					} else {
						client = wcohere.DefaultClient
						if client == nil {
							mu.Lock()
							defer mu.Unlock()
							errs = append(errs, fmt.Errorf("Cohere: we did not get an apikey for this request nor is a default client initialized"))
							return
						}
					}
					resp, err := client.Embed(cohere.EmbedOptions{
						Model:    t.Model,
						Truncate: t.Truncate,
						Texts:    texts,
					})
					mu.Lock()
					defer mu.Unlock()
					if err != nil {
						errs = append(errs, fmt.Errorf("Cohere: %v", err))
						return
					}
					for i := 0; i < len(resp.Embeddings); i++ {
						ret[i].byprovider64["cohere"] = resp.Embeddings[i]
					}
				}()
			}
		}
	}

	wg.Wait()

	return ret, nil
}

func (m *Embedder) Embed(text string, opts ...WithProviderOption) (*Embedding, error) {
	embs, err := m.BatchEmbed([]string{text}, opts...)
	if err != nil {
		return nil, err
	}

	if len(embs) != 1 {
		return nil, fmt.Errorf("We caught an internal inconsistency in our code, please report it")
	}

	return embs[0], nil
}

func (em *Embedding) GetByProvider(provider providerIden) ([]float64, error) {
	switch provider {
	default:
		panic(fmt.Errorf("should not happen: %s", provider))
	case "openai":
		if len(em.byprovider32["openai"]) == 0 {
			return nil, fmt.Errorf("This embedding does not contain provider openai")
		}
		return Float32ToFloat64(em.byprovider32["openai"]), nil
	case "cohere":
		if len(em.byprovider64["cohere"]) == 0 {
			return nil, fmt.Errorf("This embedding does not contain provider cohere")
		}
		return em.byprovider64["cohere"], nil
	}
}

func (em *Embedding) GetByProvider32(provider providerIden) ([]float32, error) {
	switch provider {
	default:
		panic(fmt.Errorf("should not happen: %s", provider))
	case "openai":
		if len(em.byprovider32["openai"]) == 0 {
			return nil, fmt.Errorf("This embedding does not contain provider openai")
		}
		return em.byprovider32["openai"], nil
	case "cohere":
		if len(em.byprovider64["cohere"]) == 0 {
			return nil, fmt.Errorf("This embedding does not contain provider cohere")
		}
		return Float64ToFloat32(em.byprovider64["cohere"]), nil
	}
}

func (em *Embedding) Get32() ([]float32, error) {
	var ret []float32

	if len(em.byprovider32)+len(em.byprovider64) > 1 {
		return nil, fmt.Errorf("This embedding contains multiple ones for different providers, use the GetByProvider or GetByProvider32 methods instead")
	}

	if len(em.byprovider32) > 0 {
		for _, p := range em.byprovider32 {
			return p, nil
		}
	}

	for _, p := range em.byprovider64 {
		return Float64ToFloat32(p), nil
	}

	return ret, nil
}

func (em *Embedding) Get() ([]float64, error) {
	if len(em.byprovider32)+len(em.byprovider64) > 1 {
		return nil, fmt.Errorf("This embedding contains multiple ones for different providers, use the GetByProvider or GetByProvider32 methods instead")
	}

	if len(em.byprovider32) > 0 {
		for _, p := range em.byprovider32 {
			return Float32ToFloat64(p), nil
		}
	}

	for _, p := range em.byprovider64 {
		return p, nil
	}

	panic("")
}

func (em *Embedding) Set(vector []float64) error {
	if len(em.byprovider32)+len(em.byprovider64) > 1 {
		return fmt.Errorf("This embedding contains multiple providers, use the SetByProvider or SetByProvider32 methods instead")
	}

	if len(em.byprovider32) > 0 {
		for p := range em.byprovider32 {
			em.byprovider32[p] = Float64ToFloat32(vector)
			return nil
		}
	}

	for p := range em.byprovider64 {
		em.byprovider64[p] = vector
		return nil
	}

	panic("")
}

func (em *Embedding) Set32(vector []float32) error {
	if len(em.byprovider32)+len(em.byprovider64) > 1 {
		return fmt.Errorf("This embedding contains multiple providers, use the SetByProvider or SetByProvider32 methods instead")
	}

	if len(em.byprovider32) > 0 {
		for p := range em.byprovider32 {
			em.byprovider32[p] = vector
			return nil
		}
	}

	for p := range em.byprovider64 {
		em.byprovider64[p] = Float32ToFloat64(vector)
		return nil
	}

	panic("")
}

func (emb *Embedding) SetByProvider(provider providerIden, vector []float64) error {
	switch provider {
	default:
		panic(fmt.Errorf("should not happen: %s", provider))
	case "openai":
		emb.byprovider32["openai"] = Float64ToFloat32(vector)
	case "cohere":
		emb.byprovider64["cohere"] = vector
	}

	return nil
}

func (emb *Embedding) SetByProvider32(provider providerIden, vector []float32) error {
	switch provider {
	default:
		panic(fmt.Errorf("should not happen: %s", provider))
	case "openai":
		emb.byprovider32["openai"] = vector
	case "cohere":
		emb.byprovider64["cohere"] = Float32ToFloat64(vector)
	}

	return nil
}
