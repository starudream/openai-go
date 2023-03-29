package openai

import (
	"net/http"
)

type CompletionsReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type CompletionsResp struct {
	Id      string               `json:"id"`
	Object  string               `json:"object"`
	Created int                  `json:"created"`
	Model   string               `json:"model"`
	Choices []*CompletionsChoice `json:"choices"`
	Usage   *CompletionsUsage    `json:"usage"`
}

type CompletionsChoice struct {
	Text         string              `json:"text"`
	Message      *CompletionsMessage `json:"message"`
	Index        int                 `json:"index"`
	FinishReason string              `json:"finish_reason"`
}

func (p *CompletionsResp) GetText() []string {
	if p == nil || len(p.Choices) == 0 {
		return nil
	}
	vs := make([]string, len(p.Choices))
	for i, choice := range p.Choices {
		if choice != nil {
			vs[i] = choice.Text
		}
	}
	return vs
}

func (p *CompletionsResp) GetMessage() []CompletionsMessage {
	if p == nil || len(p.Choices) == 0 {
		return nil
	}
	vs := make([]CompletionsMessage, len(p.Choices))
	for i, choice := range p.Choices {
		if choice != nil && choice.Message != nil {
			vs[i] = *choice.Message
		}
	}
	return vs
}

type CompletionsUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionsMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *Client) Completions(req CompletionsReq) (CompletionsResp, *http.Response, error) {
	res := CompletionsResp{}
	resp, err := c.Exec(http.MethodPost, "completions", req, &res)
	if err != nil {
		return res, nil, err
	}
	return res, resp.RawResponse, nil
}
