package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (c *Client) Pre(method, path string, opt any) (string, map[string]string, any, error) {
	u := *c.baseURL
	p, err := url.PathUnescape(path)
	if err != nil {
		return "", nil, nil, err
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
			qvs, qe := query.Values(opt)
			if qe != nil {
				return "", nil, nil, qe
			}
			u.RawQuery = qvs.Encode()
		}
	case http.MethodPost:
		headers["Content-Type"] = "application/json"
		if opt != nil {
			body = opt
		}
	}

	return u.String(), headers, body, nil
}

func (c *Client) Exec(method, path string, opt, res any) (*resty.Response, error) {
	u, hds, body, err := c.Pre(method, path, opt)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.
		NewRequest().
		SetHeaders(hds).
		SetBody(body).
		SetResult(res).
		SetError(&BaseResp{}).
		Execute(method, u)
	if err != nil {
		return nil, err
	}

	if v, ok := resp.Error().(*BaseResp); ok {
		return resp, fmt.Errorf("status: %d, type: %s, msg: %s", resp.StatusCode(), v.Error.Type, v.Error.Message)
	}

	return resp, nil
}

type BaseResp struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func (c *Client) rawPost(path string, opt any) (*http.Response, error) {
	u, hds, body, err := c.Pre(http.MethodPost, path, opt)
	if err != nil {
		return nil, err
	}

	bs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	hds["Accept"] = "text/event-stream"
	for k, v := range hds {
		req.Header.Set(k, v)
	}

	return c.client.GetClient().Do(req)
}

func newRestyClient() *resty.Client {
	c := resty.New()
	c.SetDebug(Debug())
	c.SetTimeout(5 * time.Minute)
	c.SetDisableWarn(true)
	c.SetLogger(&restyLogger{})

	return c
}
