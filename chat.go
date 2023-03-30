package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/starudream/openai-go/internal/sse"
)

type ChatCompletionsReq struct {
	Model            string               `json:"model"`
	Messages         []CompletionsMessage `json:"messages"`
	Temperature      float64              `json:"temperature,omitempty"`
	TopP             float64              `json:"top_p,omitempty"`
	N                int                  `json:"n,omitempty"`
	Stream           bool                 `json:"stream,omitempty"`
	MaxTokens        int                  `json:"max_tokens,omitempty"`
	PresencePenalty  float64              `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64              `json:"frequency_penalty,omitempty"`
}

func (c *Client) ChatCompletions(req ChatCompletionsReq) (CompletionsResp, *http.Response, error) {
	req.Stream = false
	res := CompletionsResp{}
	resp, err := c.Exec(http.MethodPost, "chat/completions", req, &res)
	if err != nil {
		return res, nil, err
	}
	return res, resp.RawResponse, nil
}

var (
	headerData     = []byte("data: ")
	headerDataDone = []byte("data: [DONE]")
)

type ChatCompletionsChan struct {
	evtCh chan CompletionsResp
	errCh chan error
}

func (c *Client) ChatCompletionsStream(req ChatCompletionsReq) (*ChatCompletionsChan, *http.Response, error) {
	req.Stream = true
	resp, err := c.rawPost("chat/completions", req)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bs, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("unexpected status code: %d, msg: %s", resp.StatusCode, bs)
	}

	res := &ChatCompletionsChan{
		evtCh: make(chan CompletionsResp),
		errCh: make(chan error),
	}

	reader := sse.NewEventStreamReader(resp.Body, 4096)

	go func() {
		defer close(res.errCh)

		for {
			bs, re := reader.ReadEvent()
			if re != nil {
				if !errors.Is(re, io.EOF) {
					res.errCh <- fmt.Errorf("read fail: %w", re)
				}
				break
			}

			if len(bs) == 0 || bytes.HasPrefix(bs, headerDataDone) {
				break
			}

			switch {
			case bytes.HasPrefix(bs, headerData):
				v := CompletionsResp{}
				e := json.Unmarshal(bytes.TrimPrefix(bs, headerData), &v)
				if e != nil {
					res.errCh <- fmt.Errorf("unmarshal fail: %w", e)
				} else {
					res.evtCh <- v
				}
			default:
				res.errCh <- fmt.Errorf("unknown event: %s", bs)
			}
		}
	}()

	return res, resp, nil
}
