package linkding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ListTagsParams struct {
	Limit  int
	Offset int
}

type ListTagsResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Tag  `json:"results"`
}

type Tag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	DateAdded time.Time `json:"date_added"`
}

type CreateTagRequest struct {
	Name string `json:"name"`
}

func (c *Client) ListTags(params ListTagsParams) (*ListTagsResponse, error) {
	path := buildTagsQueryString("/api/tags", params)

	body, err := c.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	result := &ListTagsResponse{}
	if err := json.NewDecoder(body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetTag(id int) (*Tag, error) {
	body, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/api/tags/%d/", id), nil)
	if err != nil {
		return nil, err
	}

	tag := &Tag{}
	if err := json.NewDecoder(body).Decode(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

func (c *Client) CreateTag(name string) (*Tag, error) {
	body, err := c.makeRequest(http.MethodPost, "/api/tags/", CreateTagRequest{Name: name})
	if err != nil {
		return nil, err
	}

	tag := &Tag{}
	if err := json.NewDecoder(body).Decode(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

func buildTagsQueryString(path string, params ListTagsParams) string {
	values := url.Values{}

	if params.Limit > 0 {
		values.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.Offset > 0 {
		values.Set("offset", strconv.Itoa(params.Offset))
	}

	if len(values) > 0 {
		return fmt.Sprintf("%s?%s", path, values.Encode())
	}

	return path
}
