package openai

type ClientOptionFunc func(c *Client) error

func WithBaseURL(urlStr string) ClientOptionFunc {
	return func(c *Client) error {
		return c.setBaseURL(urlStr)
	}
}

func WithUserAgent(userAgent string) ClientOptionFunc {
	return func(c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}
