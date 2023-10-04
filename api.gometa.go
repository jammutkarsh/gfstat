package main

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
	Username  string       `json:"username"`
	Type      string       `json:"type"` // User or Organization
	HTMLURL   string       `json:"html_url"`
	Followers []MetaFollow `json:"followers"`
	Following []MetaFollow `json:"following"`
}

func (c *CurrentUser) setMetadata(username, htmlURL, accountType string) {
	c.Username = username
	c.HTMLURL = htmlURL
	c.Type = accountType
}

// FollowDiff gives a intersection of CurrentUser.Followers and CurrentUser.Following
func (c CurrentUser) FollowDiff() {

}
