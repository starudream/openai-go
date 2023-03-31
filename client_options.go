package openai

import (
	"github.com/go-resty/resty/v2"
)

type ClientOptionFunc func(c *Client) error

func WithBaseURL(urlStr string) ClientOptionFunc {
	return func(c *Client) error {
		return c.setBaseURL(urlStr)
	}
}

func WithClient(client *resty.Client) ClientOptionFunc {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

func WithUserAgent(userAgent string) ClientOptionFunc {
	return func(c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}
