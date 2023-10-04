package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// GitHubAPI is the response from the GitHub API
type GitHubAPI struct {
	Username          string    `json:"login"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// GETUserData retrieves the user data from the GitHub API
func (g *GitHubAPI) GETUserData(username string) (err error) {
	response, err := http.Get(requestURL + username)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Error: User Not Found")
	}

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

func (g GitHubAPI) GETFollowers(p *CurrentUser) error {
	response, err := http.Get(g.FollowersURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Error: Couldn't find followers")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var followersGitHubData []GitHubAPI
	if err := json.Unmarshal(body, &followersGitHubData); err != nil {
		return err
	}

	for _, v := range followersGitHubData {
		p.Followers = append(p.Followers, MetaFollow{
			Username: v.Username,
			HTMLURL:  v.HTMLURL,
		})
	}

	return nil
}

func (g GitHubAPI) GETFollowing(c *CurrentUser) error {
	response, err := http.Get(g.FollowingURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Error: Couldn't find following")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var followingGitHubData []GitHubAPI
	if err := json.Unmarshal(body, &followingGitHubData); err != nil {
		return err
	}

	for _, v := range followingGitHubData {
		c.Following = append(c.Following, MetaFollow{
			Username: v.Username,
			HTMLURL:  v.HTMLURL,
		})
	}

	return nil
}

// GetUsername returns the username.
// It is specifically tied to this struct because GitHubAPI is the first data structure
// which records all the data from API call.
func (g *GitHubAPI) GetUsername() (username string) {
	return g.Username
}
