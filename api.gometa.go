package main

// MetaFollow stands for metadata of a follower/following.
type MetaFollow struct {
	Username string `json:"login"`
	HTMLURL  string `json:"html_url"`
}

// GoProfile is more abstracted version of the UserProfile struct.
// It holds the followers and following raw data, which could be processed immediately
// instead of having to make multiple API calls.
// This struct follows created for the same reason as UserProfile,
// to simplify the data structure and is subject to change based on needs.
type GoUserProfile struct {
	CurrentUser string
	Followers   []MetaFollow
	Following   []MetaFollow
}
