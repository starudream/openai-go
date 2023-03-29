package openai

import (
	"testing"
)

func TestClient_Completions(t *testing.T) {
	v, resp, err := client.Completions(CompletionsReq{
		Model:  "text-davinci-003",
		Prompt: "what's your name?",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", GetRatelimit(resp))

	vs := v.GetText()
	for i := 0; i < len(vs); i++ {
		t.Logf("%02d: %#v", i+1, vs[i])
	}
}
