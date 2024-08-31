package linkding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		http:    &http.Client{},
	}
}

func (c *Client) makeRequest(method, endpoint string, payload interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(payloadBytes)
	}

	req, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		res.Body.Close()
		return nil, fmt.Errorf("linkding: status code=%d", res.StatusCode)
	}

	return res.Body, nil
}
