package dns_1cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultTimeout = time.Second
	apiUrl         = "https://api.1cloud.ru/dns"
)

type DNS1Cloud struct {
	apiKey  string
	timeout time.Duration
	client  *http.Client
}

func New(apiKey string, opts ...OptFunc) *DNS1Cloud {
	c := &DNS1Cloud{
		apiKey: apiKey,
	}

	for _, f := range opts {
		f(c)
	}

	if c.timeout == 0 {
		c.timeout = defaultTimeout
	}

	c.client = &http.Client{
		Timeout: c.timeout,
	}

	return c
}

type OptFunc func(*DNS1Cloud)

func WithTimeout(t time.Duration) OptFunc {
	return func(c *DNS1Cloud) {
		c.timeout = t
	}
}

type command struct {
	method   string
	endpoint string
	params   interface{}
}

func (c *DNS1Cloud) send(ctx context.Context, cmd command, response interface{}) error {
	req, err := c.getRequest(cmd)
	if err != nil {
		return errors.Wrap(err, "could not get request")
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not do http request")
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	if err = dec.Decode(response); err != nil {
		return errors.Wrap(err, "could not unmarshal response")
	}

	return nil
}

func (c *DNS1Cloud) getRequest(cmd command) (*http.Request, error) {
	url := apiUrl
	if len(cmd.endpoint) > 0 {
		url = strings.Join([]string{url, cmd.endpoint}, "/")
	}

	var (
		body io.Reader
	)

	if cmd.params != nil {
		b, err := json.Marshal(cmd.params)
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal request parameters")
		}

		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(
		cmd.method,
		url,
		body,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not make request object")
	}

	req.Header.Add("Content-Type", "application/json")
	if len(c.apiKey) > 0 {
		req.Header.Add("Authorization", "Bearer "+c.apiKey)
	}

	return req, nil
}
