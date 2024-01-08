package client

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type access struct {
	AccessToken string `json:"access_token"`
	Scope       string // Scope lets us know what rights we have to the user's account
}

var GithubPublicID = os.Getenv("GH_BASIC_CLIENT_ID") // like public key

var (
	githubServerSecret = os.Getenv("GH_BASIC_SECRET_ID") // like private key
	ctx                = context.Background()
)

func init() {
	// check for env vars
	if GithubPublicID == "" || githubServerSecret == "" {
		log.Fatal("GH_BASIC_CLIENT_ID and GH_BASIC_SECRET_ID must be set")
		os.Exit(1)
	}
}

// getAccessToken returns the access token from the GitHub OAuth2.0 API
func GetAccessToken(w http.ResponseWriter, r *http.Request) (creds access) {
	sessionToken := r.URL.Query().Get("code")
	body := url.Values{"client_id": {GithubPublicID}, "client_secret": {githubServerSecret}, "code": {sessionToken}, "accept": {"json"}}

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
func GetGitHubClient(accessToken *string) *github.Client {
	tokenString := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *accessToken},
	)
	client  := oauth2.NewClient(ctx, tokenString)
	return github.NewClient(client)
}

// getGitHubUser returns the GitHub user associated with the provided GitHub client
func GetGitHubUser(client *github.Client) (*github.User, *github.Response) {
	user, resp , err := client.Users.Get(ctx, "")
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	return user, resp
}
