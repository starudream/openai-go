package openai

import (
	"net/http"
)

type ListModelResp struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

type Model struct {
	Id         string            `json:"id"`
	Object     string            `json:"object"`
	Created    int               `json:"created"`
	OwnedBy    string            `json:"owned_by"`
	Permission []ModelPermission `json:"permission"`
	Root       string            `json:"root"`
}

type ModelPermission struct {
	Id                 string `json:"id"`
	Object             string `json:"object"`
	Created            int    `json:"created"`
	AllowCreateEngine  bool   `json:"allow_create_engine"`
	AllowSampling      bool   `json:"allow_sampling"`
	AllowLogprobs      bool   `json:"allow_logprobs"`
	AllowSearchIndices bool   `json:"allow_search_indices"`
	AllowView          bool   `json:"allow_view"`
	AllowFineTuning    bool   `json:"allow_fine_tuning"`
	Organization       string `json:"organization"`
	Group              any    `json:"group"`
}

func (c *Client) ListModels() ([]Model, *http.Response, error) {
	res := ListModelResp{}
	resp, err := c.Exec(http.MethodGet, "models", nil, &res)
	if err != nil {
		return nil, nil, err
	}
	return res.Data, resp.RawResponse, nil
}

func (c *Client) RetrieveModel(id string) (Model, *http.Response, error) {
	res := Model{}
	resp, err := c.Exec(http.MethodGet, "models/"+id, nil, &res)
	if err != nil {
		return res, nil, err
	}
	return res, resp.RawResponse, nil
}
