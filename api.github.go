package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	port         = "3639"
	accessCTXKey = "access_token"
)

// Using Access Token and GitHub SDK can facilitate the use of GitHub API directly to structs.
type BasicPageData struct {
	User           github.User
	Mutuals        []MetaFollow
	IDontFollow    []MetaFollow
	TheyDontFollow []MetaFollow
}

// IndexPageData is the data for the index page template
type IndexPageData struct {
	ClientID string
}

// access is the response from the GitHub OAuth2.0 API
type access struct {
	AccessToken string `json:"access_token"`
	Scope       string // Scope lets us know what rights we have to the user's account
}

var (
	// GitHub OAuth Config
	githubPublicID     = os.Getenv("GH_BASIC_CLIENT_ID") // like public key
	githubServerSecret = os.Getenv("GH_BASIC_SECRET_ID") // like private key
	// Frontend
	indexPageData = IndexPageData{githubPublicID}
	// Context
	internalGitHubCtx = context.Background()
)

func serveWebApp() {
	fmt.Println("http://127.0.0.1:3639")
	http.HandleFunc("/", Index)
	http.HandleFunc("/result", Result)
	http.NotFoundHandler()
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

// The Index function renders the index page template and sends it as a response to the client.
func Index(w http.ResponseWriter, r *http.Request) {
	indexPage := template.Must(template.New("index.tmpl").ParseFiles("views/index.tmpl"))
	if err := indexPage.Execute(w, indexPageData); err != nil {
		log.Println(err)
	}
}

// The Result function renders the result page template and sends it as a response to the client.
func Result(w http.ResponseWriter, r *http.Request) {
	// I need to abstract my GitHUB OAuth2.0 API call to a function
	// using the token, I need to display the result.
	// But Using HTMX, I can display the result on the same page.
	// Need to figure out what happens when I make the make the callback to the same page.
	accessKeys := getAccessToken(w, r)

	client := getGitHubClient(&accessKeys.AccessToken)

	user := getGitHubUser(client)

	basicPageData := BasicPageData{*user, nil, nil, nil}

	render := template.Must(template.New("basic.tmpl").ParseFiles("views/basic.tmpl"))
	if err := render.Execute(w, basicPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
