package uni

import (
	"fmt"
	"math"
	"sort"
	"strings"
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

func NewEmbedding() *Embedding {
	return &Embedding{
		byprovider32: map[string][]float32{},
		byprovider64: map[string][]float64{},
	}
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

func (emb *Embedding) Range(fn func(provider WithProviderOption, vec []float64) error) error {
	for pname, vec := range emb.byprovider32 {
		err := fn(providerIden(pname), Float32ToFloat64(vec))
		if err != nil {
			return err
		}
	}

	for pname, vec := range emb.byprovider64 {
		err := fn(providerIden(pname), vec)
		if err != nil {
			return err
		}
	}

	return nil
}

func (emb *Embedding) Range32(fn func(provider WithProviderOption, vec []float32) error) error {
	for pname, vec := range emb.byprovider32 {
		err := fn(providerIden(pname), vec)
		if err != nil {
			return err
		}
	}

	for pname, vec := range emb.byprovider64 {
		err := fn(providerIden(pname), Float64ToFloat32(vec))
		if err != nil {
			return err
		}
	}

	return nil
}

func (emb *Embedding) getprovidershash() string {
	var providers []string

	for p := range emb.byprovider32 {
		providers = append(providers, p)
	}

	for p := range emb.byprovider64 {
		providers = append(providers, p)
	}

	sort.Strings(providers)

	return strings.Join(providers, "-")
}

// embeddings must all have the same providers
func GetMinMaxConcatenatedEmbedding(embeddings []*Embedding) (*Embedding, error) {
	ret := NewEmbedding()

	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings provided")
	}

	var phash string
	for i := 0; i < len(embeddings); i++ {
		if phash == "" {
			phash = embeddings[i].getprovidershash()
		} else {
			tmp := embeddings[i].getprovidershash()
			if tmp != phash {
				return nil, fmt.Errorf("All of the embeddings must have exactly the same providers")
			}
		}
	}

	for pname, p := range embeddings[0].byprovider32 {
		d := len(p)
		minVec := make([]float32, d)
		maxVec := make([]float32, d)

		for i := 0; i < d; i++ {
			minVal := float32(math.MaxFloat32)
			maxVal := float32(-math.MaxFloat32)

			for _, emb := range embeddings {
				if emb.byprovider32[pname][i] < minVal {
					minVal = emb.byprovider32[pname][i]
				}
				if emb.byprovider32[pname][i] > maxVal {
					maxVal = emb.byprovider32[pname][i]
				}
			}

			minVec[i] = minVal
			maxVec[i] = maxVal
		}

		concat := make([]float32, 0, len(minVec)+len(maxVec))
		concat = append(concat, minVec...)
		concat = append(concat, maxVec...)

		ret.byprovider32[pname] = concat
	}

	for pname, p := range embeddings[0].byprovider64 {
		d := len(p)
		minVec := make([]float64, d)
		maxVec := make([]float64, d)

		for i := 0; i < d; i++ {
			minVal := math.MaxFloat64
			maxVal := -math.MaxFloat64

			for _, emb := range embeddings {
				if emb.byprovider64[pname][i] < minVal {
					minVal = emb.byprovider64[pname][i]
				}
				if emb.byprovider64[pname][i] > maxVal {
					maxVal = emb.byprovider64[pname][i]
				}
			}

			minVec[i] = minVal
			maxVec[i] = maxVal
		}

		concat := make([]float64, 0, len(minVec)+len(maxVec))
		concat = append(concat, minVec...)
		concat = append(concat, maxVec...)

		ret.byprovider64[pname] = concat
	}

	return ret, nil
}
