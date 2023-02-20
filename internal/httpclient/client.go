package httpclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

var ErrStatusCode = errors.New("response status code")

type Client struct {
	httpClient *http.Client
}

func NewClient(cfg *Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				DialContext:         (&net.Dialer{Timeout: cfg.DialerTimeout}).DialContext,
				TLSHandshakeTimeout: cfg.TLSHandshakeTimeout,
			},
		},
	}
}

// DoRequest does request and return request body bytes.
func (c *Client) DoRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(ioutil.Discard, resp.Body)

		return nil, fmt.Errorf("%w do request: %v", ErrStatusCode, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return bodyBytes, nil
}
