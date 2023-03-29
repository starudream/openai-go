package openai

import (
	"testing"
)

func TestClient_ListModels(t *testing.T) {
	v, _, err := client.ListModels()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", v)
}

func TestClient_RetrieveModel(t *testing.T) {
	v, _, err := client.RetrieveModel("gpt-3.5-turbo")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", v)
}
