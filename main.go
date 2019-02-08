// Package dns1cloud implements methods for API of 1Cloud's DNS hosting
package dns1cloud

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
	defaultAPIHost = "https://api.1cloud.ru"
)

// DNS1Cloud represents a client for API of 1Cloud's DNS hosting
type DNS1Cloud struct {
	apiHost string
	apiKey  string
	timeout time.Duration
	client  *http.Client
}

// New creates and return new DNS1Cloud
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

	if len(c.apiHost) == 0 {
		c.apiHost = defaultAPIHost
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

func WithApiHost(apiHost string) OptFunc {
	return func(c *DNS1Cloud) {
		c.apiHost = apiHost
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
	url := c.apiHost
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
