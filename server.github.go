package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

// access is the response from the GitHub OAuth2.0 API
type access struct {
	AccessToken string `json:"access_token"`
	Scope       string // Scope lets us know what rights we have to the user's account
}

var (
	// GitHub OAuth Config
	githubPublicID     = os.Getenv("GH_BASIC_CLIENT_ID") // like public key
	githubServerSecret = os.Getenv("GH_BASIC_SECRET_ID") // like private key
	// Context
	internalGitHubCtx = context.Background()
)

func init() {
	// check for env vars
	if githubPublicID == "" || githubServerSecret == "" {
		log.Fatal("GH_BASIC_CLIENT_ID and GH_BASIC_SECRET_ID must be set")
		os.Exit(1)
	}
}

// getAccessToken returns the access token from the GitHub OAuth2.0 API
func getAccessToken(w http.ResponseWriter, r *http.Request) (creds access) {
	sessionToken := r.URL.Query().Get("code")
	body := url.Values{"client_id": {githubPublicID}, "client_secret": {githubServerSecret}, "code": {sessionToken}, "accept": {"json"}}

	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(body.Encode()))
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return access{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Unable to Authenticate With GitHub", http.StatusUnauthorized)
		return access{}
	}

	if err := json.NewDecoder(resp.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return access{}
	}
	return creds
}

// Authenticates GitHub Client with provided OAuth access token
// this client allows us to make make changes directly to the user's GitHub account
// without needing to manually enter various URLs and tokens
func getGitHubClient(accessToken *string) *github.Client {
	ctx := internalGitHubCtx
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// getGitHubUser returns the GitHub user associated with the provided GitHub client
func getGitHubUser(client *github.Client) *github.User {
	user, _, err := client.Users.Get(internalGitHubCtx, "")
	if err != nil {
		log.Println(err)
		return nil
	}
	return user
}

func followOverflow(user github.User) error {
	count := *user.Followers + *user.Following
	/*
		Logic: GitHub Provides a limit of 5000 Re per hour. Since this app is stateless, we don't store the user's followers and following.
		user.FollowingURL and user.FollowersURL are paginated. It only provides 100 array objects(following or followers) per request.
		Since we don't store the data, we need to make multiple requests to get the all the followers and following.
		We calculate the total number of requests required to get all the followers and following.
		We then check if the total number of requests exceeds the GitHub limit of 5000.
		totalRequests = (followers + following) / 100 && the totoalRequests > 5000, we return an error.
		Why '5000-2' ? The first request is made to get AccessToken and the second request is made to get the user.
	*/
	if count/100 > 5000-2 {
		return fmt.Errorf("you have too many followers and following (%d), it exceeds the GitHub limit of 5000", count)
	}
	return nil
}
