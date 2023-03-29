package openai

import (
	"os"
	"testing"
)

var client *Client

func TestMain(m *testing.M) {
	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		panic("OPENAI_API_KEY is not set")
	}

	cli, err := NewClient(token)
	if err != nil {
		panic(err)
	}

	client = cli

	os.Exit(m.Run())
}
