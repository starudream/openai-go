package openai

import (
	"net/http"
)

type CompletionsReq struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	Suffix           string  `json:"suffix,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
	Temperature      float64 `json:"temperature,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	N                int     `json:"n,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
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
	Text         *string             `json:"text"`
	Message      *CompletionsMessage `json:"message"`
	Delta        *CompletionsMessage `json:"delta"`
	Index        int                 `json:"index"`
	FinishReason string              `json:"finish_reason"`
}

func (p *CompletionsResp) GetMessage() []CompletionsMessage {
	if p == nil || len(p.Choices) == 0 {
		return nil
	}
	vs := make([]CompletionsMessage, len(p.Choices))
	for i, choice := range p.Choices {
		if choice.Text != nil {
			vs[i] = NewCompletionsMessage("", *choice.Text)
		} else if choice.Message != nil {
			vs[i] = *choice.Message
		} else if choice.Delta != nil {
			vs[i] = *choice.Delta
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

func NewCompletionsMessage(role, content string) CompletionsMessage {
	return CompletionsMessage{Role: role, Content: content}
}

func (c *Client) Completions(req CompletionsReq) (CompletionsResp, *http.Response, error) {
	res := CompletionsResp{}
	resp, err := c.Exec(http.MethodPost, "completions", req, &res)
	if err != nil {
		return res, nil, err
	}
	return res, resp.RawResponse, nil
}
