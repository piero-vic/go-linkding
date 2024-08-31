package linkding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ListBookmarksParams struct {
	Query  string
	Limit  int
	Offset int
	Unread bool
}

type ListBookmarksResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Bookmark `json:"results"`
}

type Bookmark struct {
	Id                    int       `json:"id"`
	URL                   string    `json:"url"`
	UserTitle             string    `json:"title"`
	UserDescription       string    `json:"description"`
	Notes                 string    `json:"notes"`
	WebsiteTitle          string    `json:"website_title"`
	WebsiteDescription    string    `json:"website_description"`
	WebArchiveSnapshotURL string    `json:"web_archive_snapshot_url"`
	FaviconURL            string    `json:"favicon_url"`
	PreviewImageURL       string    `json:"preview_image_url"`
	IsArchived            bool      `json:"is_archived"`
	Unread                bool      `json:"unread"`
	Shared                bool      `json:"shared"`
	TagNames              []string  `json:"tag_names"`
	DateAdded             time.Time `json:"date_added"`
	DateModified          time.Time `json:"date_modified"`
}

type CreateBookmarkRequest struct {
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Notes       string   `json:"notes"`
	IsArchived  bool     `json:"is_archived"`
	Unread      bool     `json:"unread"`
	Shared      bool     `json:"shared"`
	TagNames    []string `json:"tag_names"`
}

func (c *Client) ListBookmarks(params ListBookmarksParams) (*ListBookmarksResponse, error) {
	path := buildBookmarksQueryString("/api/bookmarks", params)

	body, err := c.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	result := &ListBookmarksResponse{}
	if err := json.NewDecoder(body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetBookmark(id int) (*Bookmark, error) {
	body, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/api/bookmarks/%d/", id), nil)
	if err != nil {
		return nil, err
	}

	bookmark := &Bookmark{}
	if err := json.NewDecoder(body).Decode(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

func (c *Client) CreateBookmark(payload CreateBookmarkRequest) (*Bookmark, error) {
	body, err := c.makeRequest(http.MethodPost, "/api/bookmarks/", payload)
	if err != nil {
		return nil, err
	}

	bookmark := &Bookmark{}
	if err := json.NewDecoder(body).Decode(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

func (c *Client) UpdateBookmark(id int, payload CreateBookmarkRequest) (*Bookmark, error) {
	body, err := c.makeRequest(http.MethodPut, fmt.Sprintf("/api/bookmarks/%d/", id), payload)
	if err != nil {
		return nil, err
	}

	bookmark := &Bookmark{}
	if err := json.NewDecoder(body).Decode(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

func (c *Client) DeleteBookmark(id int) error {
	_, err := c.makeRequest(http.MethodDelete, fmt.Sprintf("/api/bookmarks/%d/", id), nil)

	return err
}

func buildBookmarksQueryString(path string, params ListBookmarksParams) string {
	values := url.Values{}

	if params.Query != "" {
		values.Set("q", params.Query)
	}

	if params.Limit > 0 {
		values.Set("limit", strconv.Itoa(params.Limit))
	}

	if params.Offset > 0 {
		values.Set("offset", strconv.Itoa(params.Offset))
	}

	if params.Unread {
		values.Set("unread", "yes")
	}

	if len(values) > 0 {
		return fmt.Sprintf("%s?%s", path, values.Encode())
	}

	return path
}
