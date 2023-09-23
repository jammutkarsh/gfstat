package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type AccountType int

const (
	User AccountType = iota
	Organization
)

// UserProfile is the fundamental data structure which hold the necesary to be processed.
// It abstracts the GitHubAPI struct of unnecessary information, simplifying the data structure.
// As the app grows, this struct can be further expanded to hold more information.
// Currently it is designed to hold the bare minimum information which is required to get the follow diff.
// Potential future additions could be used to hold bio, avatar, etc. which could be used to display a profile in the frontend.
type UserProfile struct {
	Username     string      `json:"username"`
	Type         AccountType `json:"type"` // User or Organization
	GitHubURL    string      `json:"html_url"`
	FollowingURL string      `json:"following_url"`
	FollowersURL string      `json:"followers_url"`
}

// GETFollowers retrieves the followers data for a user.
func (out UserProfile) GETFollowers(p *GoUserProfile) error {
	response, err := http.Get(out.FollowersURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var g []GitHubAPI
	if err := json.Unmarshal(body, &g); err != nil {
		return err
	}

	for _, v := range g {
		p.Followers = append(p.Followers, MetaFollow{
			Username: v.Username,
			HTMLURL:  v.HTMLURL,
		})
	}

	return nil
}

// GETFollowing retrieves the followers data for a user.
func (out UserProfile) GETFollowing(p *GoUserProfile) error {
	response, err := http.Get(out.FollowingURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var g []GitHubAPI
	if err := json.Unmarshal(body, &g); err != nil {
		return err
	}

	for _, v := range g {
		p.Following = append(p.Following, MetaFollow{
			Username: v.Username,
			HTMLURL:  v.HTMLURL,
		})
	}

	return nil
}
