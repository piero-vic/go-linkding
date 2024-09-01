package linkding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ListBookmarksParams defines the parameters used when listing bookmarks.
type ListBookmarksParams struct {
	// The search query to filter bookmarks.
	Query string
	// The maximum number of bookmarks to return.
	Limit int
	// The offset for pagination.
	Offset int
	// Filter to include only unread bookmarks.
	Unread bool
}

// ListBookmarksResponse represents the response from the Linkding API when
// listing bookmarks.
type ListBookmarksResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Bookmark `json:"results"`
}

// Bookmark represents a bookmark object in the Linkding API.
type Bookmark struct {
	ID                    int       `json:"id"`
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

// CreateBookmarkRequest represents the request body when creating or updating
// bookmarks.
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

// ListBookmarks retrieves a list of bookmarks from Linkding based on the
// provided parameters.
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

// GetBookmark retrieves a single bookmark from Linkding.
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

// CreateBookmark creates a new bookmark in Linkding using the provided payload.
//
// Warning: Ensure that the TagNames property in the CreateBookmarkRequest is
// initialized (even if empty) to avoid nil pointer issues.
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

// UpdateBookmark updates an existing bookmark in Linkding using the provided
// payload.
//
// Warning: Ensure that the TagNames property in the CreateBookmarkRequest is
// initialized (even if empty) to avoid nil pointer issues.
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

// DeleteBookmark deletes a bookmark from Linkding.
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
