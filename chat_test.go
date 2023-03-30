package openai

import (
	"testing"
)

func TestClient_ChatCompletions(t *testing.T) {
	v, resp, err := client.ChatCompletions(ChatCompletionsReq{
		Model: modelGPT35,
		Messages: []CompletionsMessage{
			NewCompletionsMessage("user", "what's your name?"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", GetRatelimit(resp))

	vs := v.GetMessage()
	for i := 0; i < len(vs); i++ {
		t.Logf("%02d: %s %#v", i+1, vs[i].Role, vs[i].Content)
	}
}

func TestClient_ChatCompletionsStream(t *testing.T) {
	res, _, err := client.ChatCompletionsStream(ChatCompletionsReq{
		Model: modelGPT35,
		Messages: []CompletionsMessage{
			NewCompletionsMessage("user", "what's your name?"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case v := <-res.evtCh:
			vs := v.GetMessage()
			for i := 0; i < len(vs); i++ {
				t.Logf("%02d: %s %#v", i+1, vs[i].Role, vs[i].Content)
			}
		case e, ok := <-res.errCh:
			if e != nil {
				t.Fatal(err)
			}
			if !ok {
				res.errCh = nil
			}
		}
		if res.errCh == nil {
			break
		}
	}
}
