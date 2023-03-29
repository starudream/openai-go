package openai

import (
	"net/http"
)

type ChatCompletionsReq struct {
	Model    string               `json:"model"`
	Messages []CompletionsMessage `json:"messages"`
}

func (c *Client) ChatCompletions(req ChatCompletionsReq) (CompletionsResp, *http.Response, error) {
	res := CompletionsResp{}
	resp, err := c.Exec(http.MethodPost, "chat/completions", req, &res)
	if err != nil {
		return res, nil, err
	}
	return res, resp.RawResponse, nil
}
