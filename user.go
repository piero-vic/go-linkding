package linkding

import (
	"encoding/json"
	"net/http"
)

type UserPreferences struct {
	Theme                 string `json:"theme"`
	BookmarkDateDisplay   string `json:"bookmark_date_display"`
	BookmarkLinkTarget    string `json:"bookmark_link_target"`
	WebArchiveIntegration string `json:"web_archive_integration"`
	TagSearch             string `json:"tag_search"`
	EnableSharing         bool   `json:"enable_sharing"`
	EnablePublicSharing   bool   `json:"enable_public_sharing"`
	EnableFavicons        bool   `json:"enable_favicons"`
	DisplayURL            bool   `json:"display_url"`
	PermanentNotes        bool   `json:"permanent_notes"`
	SearchPreferences     struct {
		Sort   string `json:"sort"`
		Shared string `json:"shared"`
		Unread string `json:"unread"`
	} `json:"search_preferences"`
}

func (c *Client) GetUserPreferences() (*UserPreferences, error) {
	body, err := c.makeRequest(http.MethodGet, "/api/user/profile/", nil)
	if err != nil {
		return nil, err
	}

	userPreferences := &UserPreferences{}
	if err := json.NewDecoder(body).Decode(userPreferences); err != nil {
		return nil, err
	}

	return userPreferences, nil
}