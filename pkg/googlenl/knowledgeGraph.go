package googlenl

import (
	"encoding/json"
)

type SearchResult struct {
	Type        string             `json:"@type"`
	Result      SearchResultEntity `json:"result"`
	ResultScore float64            `json:"resultScore"`
}

type SearchResultEntity struct {
	Id                  string              `json:"@id"`
	Type                []string            `json:"@type"`
	Name                string              `json:"name"`
	Image               EntityImage         `json:"image"`
	DetailedDescription DetailedDescription `json:"detailedDescription"`
}

type EntityImage struct {
	ContentUrl string `json:"contentUrl"`
	Url        string `json:"url"`
}

type DetailedDescription struct {
	ArticleBody string `json:"articleBody"`
	License     string `json:"license"`
	Url         string `json:"url"`
}

type KnowledgeGraphSearchResponse struct {
	ItemListElement []*SearchResult `json:"itemListElement,omitempty"`
}

func SearchKnowledgeGraph(query string) (*KnowledgeGraphSearchResponse, error) {
	req := KnowledgeGraphClient.Entities.Search()
	req.Query(query)

	r, err := req.Do()
	if err != nil {
		return nil, err
	}

	resp := &KnowledgeGraphSearchResponse{}

	for i := 0; i < len(r.ItemListElement); i++ {
		b, err := json.Marshal(r.ItemListElement[i])
		if err != nil {
			return nil, err
		}
		sr := &SearchResult{}
		err = json.Unmarshal(b, sr)
		if err != nil {
			return nil, err
		}
		resp.ItemListElement = append(resp.ItemListElement, sr)
	}

	return resp, nil
}

func SearchKnowledgeGraphByIds(ids ...string) (*KnowledgeGraphSearchResponse, error) {
	req := KnowledgeGraphClient.Entities.Search()
	req.Ids(ids...)

	r, err := req.Do()
	if err != nil {
		return nil, err
	}

	resp := &KnowledgeGraphSearchResponse{}

	for i := 0; i < len(r.ItemListElement); i++ {
		b, err := json.Marshal(r.ItemListElement[i])
		if err != nil {
			return nil, err
		}
		sr := &SearchResult{}
		err = json.Unmarshal(b, sr)
		if err != nil {
			return nil, err
		}
		resp.ItemListElement = append(resp.ItemListElement, sr)
	}

	return resp, nil
}
