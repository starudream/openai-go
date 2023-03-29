package openai

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/starudream/openai-go/internal/query"
)

const (
	DefaultBaseURL   = "https://api.openai.com/"
	DefaultVersion   = "v1/"
	DefaultUserAgent = "starudream-openai-go"
)

type Client struct {
	client *resty.Client

	baseURL *url.URL

	userAgent string

	token string
}

func NewClient(token string, options ...ClientOptionFunc) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	c := &Client{
		client:    newRestyClient(),
		userAgent: DefaultUserAgent,
		token:     token,
	}

	_ = c.setBaseURL(DefaultBaseURL)

	for i := 0; i < len(options); i++ {
		option := options[i]
		if option == nil {
			continue
		}
		if err := option(c); err != nil {
			return nil, err
		}
	}

	c.client.SetHeader("User-Agent", c.userAgent)

	return c, nil
}

func (c *Client) setBaseURL(urlStr string) error {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(baseURL.Path, DefaultVersion) {
		baseURL.Path += DefaultVersion
	}

	c.baseURL = baseURL

	return nil
}

func (c *Client) Exec(method, path string, opt, res any) (*resty.Response, error) {
	u := *c.baseURL
	p, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + p

	headers := map[string]string{
		"Authorization": "Bearer " + c.token,
	}

	var body any

	switch method {
	case http.MethodGet:
		if opt != nil {
			qvs, err := query.Values(opt)
			if err != nil {
				return nil, err
			}
			u.RawQuery = qvs.Encode()
		}
	case http.MethodPost:
		headers["Content-Type"] = "application/json"
		if opt != nil {
			body = opt
		}
	}

	resp, err := c.client.
		NewRequest().
		SetHeaders(headers).
		SetBody(body).
		SetResult(res).
		Execute(method, u.String())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func newRestyClient() *resty.Client {
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	c := resty.New()
	c.SetDebug(debug)
	c.SetTimeout(5 * time.Minute)
	c.SetDisableWarn(true)
	c.SetLogger(&restyLogger{})

	return c
}
