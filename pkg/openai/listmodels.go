package openai

const urlSuffix_listmodels = "v1/models"

type ListModelsRequest struct {
	// Only required if no default api key was initialized
	APIKEY string `json:"-"`
}

type ListModelsObject struct {
	Object     string `json:"object"`
	ID         string `json:"id"`
	Created    int    `json:"created"`
	OwnedBy    string `json:"owned_by"`
	Permission []struct {
		ID                 string      `json:"id"`
		Object             string      `json:"object"`
		Created            int         `json:"created"`
		AllowCreateEngine  bool        `json:"allow_create_engine"`
		AllowSampling      bool        `json:"allow_sampling"`
		AllowLogprobs      bool        `json:"allow_logprobs"`
		AllowSearchIndices bool        `json:"allow_search_indices"`
		AllowView          bool        `json:"allow_view"`
		AllowFineTuning    bool        `json:"allow_fine_tuning"`
		Organization       string      `json:"organization"`
		Group              interface{} `json:"group"`
		IsBlocking         bool        `json:"is_blocking"`
	} `json:"permission"`
	Root   string      `json:"root"`
	Parent interface{} `json:"parent"`
}

type ListModelsResponse struct {
	Data   []ListModelsObject `json:"data"`
	Object string             `json:"object"`
}

func ListModels(req *ListModelsRequest) (*ListModelsResponse, error) {
	resp := &ListModelsResponse{}

	err := request("GET", urlSuffix_listmodels, nil, resp, req.APIKEY)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
