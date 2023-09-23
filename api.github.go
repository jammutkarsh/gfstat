package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// GitHubAPI is the response from the GitHub API
type GitHubAPI struct {
	Username          string `json:"login"`
	Type              string `json:"type"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	SiteAdmin         bool   `json:"site_admin"`
}

// GETUserData retrieves the user data from the GitHub API
func (g *GitHubAPI) GETUserData(username string) (err error) {
	response, err := http.Get(requestURL + username)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &g); err != nil {
		return err
	}

	// Remove the trailing "{/other_user}" from the URL
	g.FollowingURL = g.FollowingURL[:len(g.FollowingURL)-len("{/other_user}")]

	return nil
}

// GetUsername returns the username.
// It is specifically tied to this struct because GitHubAPI is the first data structure
// which records all the data from API call.
func (g *GitHubAPI) GetUsername() (username string) {
	return g.Username
}
