package main

import "sort"

type AccountType int

const (
	User AccountType = iota
	Organization
)

// MetaFollow stands for metadata of a follower/following.
type MetaFollow struct {
	Username string `json:"username"`
	HTMLURL  string `json:"html_url"`
}

// GoProfile is more abstracted version of the UserProfile struct.
// It holds the followers and following raw data, which could be processed immediately
// instead of having to make multiple API calls.
// This struct follows created for the same reason as UserProfile,
// to simplify the data structure and is subject to change based on needs.
type CurrentUser struct {
	Username       string       `json:"username"`
	Type           string       `json:"type"` // User or Organization
	HTMLURL        string       `json:"html_url"`
	Followers      []MetaFollow `json:"followers"`
	FollowersCount int
	Following      []MetaFollow `json:"following"`
	FollowingCount int
}

func (c *CurrentUser) setMetadata(username, htmlURL, accountType string) {
	c.Username = username
	c.HTMLURL = htmlURL
	c.Type = accountType
}

// Mutuals gives the list of mutuals between followers and following.
func (c CurrentUser) Mutuals() []MetaFollow {
	// sort following and followers based on usernames
	sort.Slice(c.Followers, func(i, j int) bool {
		return c.Followers[i].Username < c.Followers[j].Username
	})
	sort.Slice(c.Following, func(i, j int) bool {
		return c.Following[i].Username < c.Following[j].Username
	})
	// find mutuals
	var mutuals []MetaFollow
	for _, follower := range c.Followers {
		for _, following := range c.Following {
			if follower.Username == following.Username {
				mutuals = append(mutuals, follower)
			}
		}
	}
	return mutuals
}
