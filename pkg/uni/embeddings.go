package uni

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/arthurweinmann/go-ai-sdk/pkg/openai"
	"github.com/arthurweinmann/go-ai-sdk/pkg/wcohere"
	api "github.com/cohere-ai/cohere-go/v2"
	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
)

type Embedder struct {
	err       error
	providers []EmbedderOption
}

type SingleProviderEmbedder struct {
	err error
	opt EmbedderOption
}

type Embedding struct {
	errByProvider map[string]error

	byprovider32 map[string][]float32
	byprovider64 map[string][]float64
}

type SingleProviderEmbedding struct {
	v32 []float32
	v64 []float64
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

func WithCohereEmbed(model, truncate, inputType, apikeyOptional string) *withCohereOption {
	return &withCohereOption{
		APIKey:    apikeyOptional,
		Model:     model,
		Truncate:  truncate,
		InputType: inputType,
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

func NewSingleProviderEmbedder(opts ...EmbedderOption) *SingleProviderEmbedder {
	emb := &SingleProviderEmbedder{}

	if len(opts) == 0 {
		emb.err = fmt.Errorf("We need one provider of embeddings")
		return emb
	}

	var providerset bool

	for i := 0; i < len(opts); i++ {
		switch t := opts[i].(type) {
		default:
			panic(fmt.Errorf("Should not happen: %T", t))
		case *withOpenAIOption:
			emb.opt = t
			if providerset {
				emb.err = fmt.Errorf("We support only one provider of embeddings for a SingleProviderEmbedder")
				return emb
			}
			providerset = true
		case *withCohereOption:
			emb.opt = t
			if providerset {
				emb.err = fmt.Errorf("We support only one provider of embeddings for a SingleProviderEmbedder")
				return emb
			}
			providerset = true
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

func (emb *Embedding) FirstError() error {
	for pname, e := range emb.errByProvider {
		if e != nil {
			return fmt.Errorf("%s: %v", pname, e)
		}
	}
	return nil
}

func (emb *Embedding) ByProviderError(provider WithProviderOption) error {
	switch t := provider.(type) {
	default:
		panic(fmt.Errorf("Should not happen: %T", t))
	case providerIden:
		switch t {
		default:
			panic(fmt.Errorf("Should not happen: %s", t))
		case "openai", "cohere":
			return emb.errByProvider[string(t)]
		}
	}
}

func EmbeddingFrom(provider WithProviderOption, vector []float64) *Embedding {
	emb := NewEmbedding()

	switch t := provider.(type) {
	default:
		panic(fmt.Errorf("Should not happen: %T", t))
	case providerIden:
		switch t {
		default:
			panic(fmt.Errorf("Should not happen: %s", t))
		case "openai":
			emb.byprovider32["openai"] = Float64ToFloat32(vector)
		case "cohere":
			emb.byprovider64["cohere"] = vector
		}
	}

	return emb
}

func EmbeddingFrom32(provider WithProviderOption, vector []float32) *Embedding {
	emb := NewEmbedding()

	switch t := provider.(type) {
	default:
		panic(fmt.Errorf("Should not happen: %T", t))
	case providerIden:
		switch t {
		default:
			panic(fmt.Errorf("Should not happen: %s", t))
		case "openai":
			emb.byprovider32["openai"] = vector
		case "cohere":
			emb.byprovider64["cohere"] = Float32ToFloat64(vector)
		}
	}

	return emb
}

func (emb *Embedding) ToSingleProvider() (*SingleProviderEmbedding, error) {
	if len(emb.byprovider32)+len(emb.byprovider64) > 1 {
		return nil, fmt.Errorf("We cannot convert an Embedding with multiple vectors for multiple providers to a SingleProviderEmbedding")
	}

	s := NewSingleProviderEmbedding()

	for _, v := range emb.byprovider32 {
		s.v32 = v
		s.v64 = nil
		break
	}

	for _, v := range emb.byprovider64 {
		s.v32 = nil
		s.v64 = v
		break
	}

	return s, nil
}

func NewSingleProviderEmbedding() *SingleProviderEmbedding {
	return &SingleProviderEmbedding{}
}

func SingleProviderEmbeddingFrom(v []float64) *SingleProviderEmbedding {
	s := NewSingleProviderEmbedding()

	s.v64 = v

	return s
}

func SingleProviderEmbeddingFrom32(v []float32) *SingleProviderEmbedding {
	s := NewSingleProviderEmbedding()

	s.v32 = v

	return s
}

func (s *SingleProviderEmbedding) ToEmbedding(provider WithProviderOption) *Embedding {
	t, ok := provider.(providerIden)
	if !ok {
		panic(fmt.Errorf("Should not happen: %T", provider))
	}

	ret := NewEmbedding()

	switch t {
	default:
		panic(fmt.Errorf("Should not happen: %s", t))
	case "openai":
		if len(s.v32) > 0 {
			ret.byprovider32["openai"] = s.v32
		} else {
			ret.byprovider32["openai"] = Float64ToFloat32(s.v64)
		}
	case "cohere":
		if len(s.v32) > 0 {
			ret.byprovider64["cohere"] = Float32ToFloat64(s.v32)
		} else {
			ret.byprovider64["cohere"] = s.v64
		}
	}

	return ret
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
			byprovider32:  map[string][]float32{},
			byprovider64:  map[string][]float64{},
			errByProvider: map[string]error{},
		}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, prov := range m.providers {

		for k := 0; k < len(texts); k += 50 {
			l := k + 50
			var tmpbatch []string
			if l < len(texts) {
				tmpbatch = texts[k:l]
			} else {
				tmpbatch = texts[k:]
			}

			switch t := prov.(type) {
			default:
				panic(fmt.Errorf("Should not happen: %T", t))
			case *withOpenAIOption:
				if useOpenAI {
					wg.Add(1)
					go func(kindex int, batch []string) {
						defer wg.Done()
						resp, err := openai.CreateEmbedding(&openai.EmbeddingRequest{
							APIKEY: t.APIKey,
							Model:  t.Model,
							Input:  batch,
						})
						mu.Lock()
						defer mu.Unlock()
						if err != nil {
							for i := 0; i < len(resp.Data); i++ {
								ret[kindex+resp.Data[i].Index].errByProvider["openai"] = err
							}
							return
						}
						for i := 0; i < len(resp.Data); i++ {
							ret[kindex+resp.Data[i].Index].byprovider32["openai"] = resp.Data[i].Embedding
						}
					}(k, tmpbatch)
				}
			case *withCohereOption:
				if useCohere {
					wg.Add(1)
					go func(kindex int, batch []string) {
						defer wg.Done()
						var client *cohereclient.Client
						if t.APIKey != "" {
							var err error
							client, err = wcohere.NewClient(t.APIKey)
							if err != nil {
								mu.Lock()
								defer mu.Unlock()
								for i := 0; i < len(batch); i++ {
									ret[kindex+i].errByProvider["cohere"] = err
								}
								return
							}
						} else {
							client = wcohere.DefaultClient
							if client == nil {
								mu.Lock()
								defer mu.Unlock()
								for i := 0; i < len(batch); i++ {
									ret[kindex+i].errByProvider["cohere"] = fmt.Errorf("Cohere: we did not get an apikey for this request nor is a default client initialized")
								}
								return
							}
						}
						params := &cohere.EmbedRequest{
							Model: &t.Model,
							Texts: batch,
						}
						truncateOpt := cohere.EmbedRequestTruncate(t.Truncate)
						if truncateOpt != "" {
							params.Truncate = &truncateOpt
						}
						inputTypeOpt := api.EmbedInputType(t.InputType)
						if inputTypeOpt != "" {
							params.InputType = &inputTypeOpt
						}
						resp, err := client.Embed(context.Background(), params)
						mu.Lock()
						defer mu.Unlock()
						if err != nil {
							for i := 0; i < len(batch); i++ {
								ret[kindex+i].errByProvider["cohere"] = err
							}
							return
						}
						for i := 0; i < len(resp.EmbeddingsFloats.Embeddings); i++ {
							ret[kindex+i].byprovider64["cohere"] = resp.EmbeddingsFloats.Embeddings[i]
						}
					}(k, tmpbatch)
				}
			}
		}
	}

	wg.Wait()

	return ret, nil
}

func (m *SingleProviderEmbedder) BatchEmbed(texts []string, opts ...WithProviderOption) ([]*SingleProviderEmbedding, error) {
	if m.err != nil {
		return nil, m.err
	}

	ret := make([]*SingleProviderEmbedding, len(texts))
	for i := 0; i < len(ret); i++ {
		ret[i] = &SingleProviderEmbedding{}
	}

	for k := 0; k < len(texts); k += 50 {
		l := k + 50
		var tmpbatch []string
		if l < len(texts) {
			tmpbatch = texts[k:l]
		} else {
			tmpbatch = texts[k:]
		}

		switch t := m.opt.(type) {
		default:
			panic(fmt.Errorf("Should not happen: %T", t))
		case *withOpenAIOption:
			resp, err := openai.CreateEmbedding(&openai.EmbeddingRequest{
				APIKEY: t.APIKey,
				Model:  t.Model,
				Input:  tmpbatch,
			})
			if err != nil {
				return nil, err
			}
			for i := 0; i < len(resp.Data); i++ {
				ret[k+resp.Data[i].Index].v32 = resp.Data[i].Embedding
			}
		case *withCohereOption:
			var client *cohereclient.Client
			if t.APIKey != "" {
				var err error
				client, err = wcohere.NewClient(t.APIKey)
				if err != nil {
					return nil, err
				}
			} else {
				client = wcohere.DefaultClient
				if client == nil {
					return nil, fmt.Errorf("Cohere: we did not get an apikey for this request nor is a default client initialized")
				}
			}
			params := &cohere.EmbedRequest{
				Model: &t.Model,
				Texts: tmpbatch,
			}
			truncateOpt := cohere.EmbedRequestTruncate(t.Truncate)
			if truncateOpt != "" {
				params.Truncate = &truncateOpt
			}
			inputTypeOpt := api.EmbedInputType(t.InputType)
			if inputTypeOpt != "" {
				params.InputType = &inputTypeOpt
			}
			resp, err := client.Embed(context.Background(), params)
			if err != nil {
				return nil, err
			}
			for i := 0; i < len(resp.EmbeddingsFloats.Embeddings); i++ {
				ret[k+i].v64 = resp.EmbeddingsFloats.Embeddings[i]
			}
		}
	}

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

func (m *SingleProviderEmbedder) Embed(text string, opts ...WithProviderOption) (*SingleProviderEmbedding, error) {
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
		if em.errByProvider != nil && em.errByProvider["openai"] != nil {
			return nil, em.errByProvider["openai"]
		}
		if len(em.byprovider32["openai"]) == 0 {
			return nil, fmt.Errorf("This embedding does not contain provider openai")
		}
		return Float32ToFloat64(em.byprovider32["openai"]), nil
	case "cohere":
		if em.errByProvider != nil && em.errByProvider["cohere"] != nil {
			return nil, em.errByProvider["cohere"]
		}
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
		if em.errByProvider != nil && em.errByProvider["openai"] != nil {
			return nil, em.errByProvider["openai"]
		}
		if len(em.byprovider32["openai"]) == 0 {
			return nil, fmt.Errorf("This embedding does not contain provider openai")
		}
		return em.byprovider32["openai"], nil
	case "cohere":
		if em.errByProvider != nil && em.errByProvider["cohere"] != nil {
			return nil, em.errByProvider["cohere"]
		}
		if len(em.byprovider64["cohere"]) == 0 {
			return nil, fmt.Errorf("This embedding does not contain provider cohere")
		}
		return Float64ToFloat32(em.byprovider64["cohere"]), nil
	}
}

func (em *SingleProviderEmbedding) Get32() []float32 {
	if len(em.v32) > 0 {
		return em.v32
	}

	return Float64ToFloat32(em.v64)
}

func (em *SingleProviderEmbedding) Get() []float64 {
	if len(em.v64) > 0 {
		return em.v64
	}

	return Float32ToFloat64(em.v32)
}

func (em *SingleProviderEmbedding) Set(vector []float64) {
	em.v64 = vector
	em.v32 = nil
}

func (em *SingleProviderEmbedding) Set32(vector []float32) {
	em.v32 = vector
	em.v64 = nil
}

func (emb *Embedding) SetByProvider(provider providerIden, vector []float64) error {
	switch provider {
	default:
		panic(fmt.Errorf("should not happen: %s", provider))
	case "openai":
		emb.byprovider32["openai"] = Float64ToFloat32(vector)
		if emb.errByProvider != nil {
			delete(emb.errByProvider, "openai")
		}
	case "cohere":
		emb.byprovider64["cohere"] = vector
		if emb.errByProvider != nil {
			delete(emb.errByProvider, "cohere")
		}
	}

	return nil
}

func (emb *Embedding) SetByProvider32(provider providerIden, vector []float32) error {
	switch provider {
	default:
		panic(fmt.Errorf("should not happen: %s", provider))
	case "openai":
		emb.byprovider32["openai"] = vector
		if emb.errByProvider != nil {
			delete(emb.errByProvider, "openai")
		}
	case "cohere":
		emb.byprovider64["cohere"] = Float32ToFloat64(vector)
		if emb.errByProvider != nil {
			delete(emb.errByProvider, "cohere")
		}
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

// embeddings' vectors must all have the same length
func GetMinMaxConcatenatedSingleProviderEmbedding(embeddings []*SingleProviderEmbedding) (*SingleProviderEmbedding, error) {
	ret := NewSingleProviderEmbedding()

	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings provided")
	}

	var count32, count64 int
	for _, emb := range embeddings {
		if len(emb.v32) > 0 {
			count32++
		}
		if len(emb.v64) > 0 {
			count64++
		}
	}

	var d int

	if len(embeddings[0].v32) > 0 {
		d = len(embeddings[0].v32)
	} else {
		d = len(embeddings[0].v64)
	}

	if count32 > count64 {
		minVec := make([]float32, d)
		maxVec := make([]float32, d)

		for i := 0; i < d; i++ {
			minVal := float32(math.MaxFloat32)
			maxVal := float32(-math.MaxFloat32)

			for _, emb := range embeddings {
				if len(emb.v32) > 0 {
					if len(emb.v32) != d {
						return nil, fmt.Errorf("We need all the vectors from all the embeddings to be of the same length")
					}
					if emb.v32[i] < minVal {
						minVal = emb.v32[i]
					}
					if emb.v32[i] > maxVal {
						maxVal = emb.v32[i]
					}
				} else {
					if len(emb.v64) != d {
						return nil, fmt.Errorf("We need all the vectors from all the embeddings to be of the same length")
					}
					if float32(emb.v64[i]) < minVal {
						minVal = float32(emb.v64[i])
					}
					if float32(emb.v64[i]) > maxVal {
						maxVal = float32(emb.v64[i])
					}
				}
			}

			minVec[i] = minVal
			maxVec[i] = maxVal
		}

		concat := make([]float32, 0, len(minVec)+len(maxVec))
		concat = append(concat, minVec...)
		concat = append(concat, maxVec...)

		ret.v32 = concat
		ret.v64 = nil
	} else {
		minVec := make([]float64, d)
		maxVec := make([]float64, d)

		for i := 0; i < d; i++ {
			minVal := float64(math.MaxFloat64)
			maxVal := float64(-math.MaxFloat64)

			for _, emb := range embeddings {
				if len(emb.v64) > 0 {
					if len(emb.v64) != d {
						return nil, fmt.Errorf("We need all the vectors from all the embeddings to be of the same length")
					}
					if emb.v64[i] < minVal {
						minVal = emb.v64[i]
					}
					if emb.v64[i] > maxVal {
						maxVal = emb.v64[i]
					}
				} else {
					if len(emb.v32) != d {
						return nil, fmt.Errorf("We need all the vectors from all the embeddings to be of the same length")
					}
					if float64(emb.v32[i]) < minVal {
						minVal = float64(emb.v32[i])
					}
					if float64(emb.v32[i]) > maxVal {
						maxVal = float64(emb.v32[i])
					}
				}
			}

			minVec[i] = minVal
			maxVec[i] = maxVal
		}

		concat := make([]float64, 0, len(minVec)+len(maxVec))
		concat = append(concat, minVec...)
		concat = append(concat, maxVec...)

		ret.v32 = nil
		ret.v64 = concat
	}

	return ret, nil
}
