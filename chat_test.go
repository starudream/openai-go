package openai

import (
	"testing"
)

func TestClient_ChatCompletions(t *testing.T) {
	v, resp, err := client.ChatCompletions(ChatCompletionsReq{
		Model: "gpt-3.5-turbo",
		Messages: []CompletionsMessage{
			{
				Role:    "user",
				Content: "How to evaluate the programming language of Golang?",
			},
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
