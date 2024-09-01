package linkding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client handles all interactions with the Linkding API.
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

// NewClient creates a new Linkding API client using the given URL and token.
//
// The URL provided must be a complete URL. It must contain a schema and the
// domain for the API. Do not include the prefix path of the API.
// e.g. "https://linkding.example.org".
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		http:    &http.Client{},
	}
}

var (
	ErrInternalServerError = errors.New("linkding: internal server error")
	ErrUnauthorized        = errors.New("linkding: unauthorized")
	ErrNotFound            = errors.New("linkding: not found")
	ErrBadRequest          = errors.New("linkding: bad request")
)

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

	switch res.StatusCode {
	case http.StatusInternalServerError:
		res.Body.Close()
		return nil, ErrInternalServerError
	case http.StatusUnauthorized:
		res.Body.Close()
		return nil, ErrUnauthorized
	case http.StatusNotFound:
		res.Body.Close()
		return nil, ErrNotFound
	case http.StatusBadRequest:
		defer res.Body.Close()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("%w (%v)", ErrBadRequest, err)
		}

		return nil, fmt.Errorf("%w (%s)", ErrBadRequest, string(bodyBytes))
	}

	return res.Body, nil
}
