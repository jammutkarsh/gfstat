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

const port = "3639"

// Using Access Token and GitHub SDK can facilitate the use of GitHub API directly to structs.
type BasicPageData struct {
	User           *github.User
	Mutuals        []MetaFollow
	IDontFollow    []MetaFollow
	TheyDontFollow []MetaFollow
}

type Access struct {
	AccessToken string `json:"access_token"`
	Scope       string // Scope lets us know what rights we have to the user's account
}

var (
	// GitHub OAuth Config
	githubPublicID     = os.Getenv("GH_BASIC_CLIENT_ID") // like public key
	githubServerSecret = os.Getenv("GH_BASIC_SECRET_ID") // like private key
	// Frontend
	indexPageData = githubPublicID
	// Context
	background = context.Background()
)

func serveWebApp() {
	fmt.Println("http://127.0.0.1:3639")
	http.HandleFunc("/", Index)
	http.HandleFunc("/success", Success)
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

func GetAccessToken(w http.ResponseWriter, r *http.Request) string {
	code := r.URL.Query().Get("code")
	values := url.Values{"client_id": {githubPublicID}, "client_secret": {githubServerSecret}, "code": {code}, "accept": {"json"}}

	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(values.Encode()))
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Print(err)
		return ""
	}
	defer resp.Body.Close()

	var access Access

	if err := json.NewDecoder(resp.Body).Decode(&access); err != nil {
		log.Println("JSON-Decode-Problem: ", err)
		return ""
	}

	return access.AccessToken
}

// Authenticates GitHub Client with provided OAuth access token
func getGitHubClient(accessToken string) *github.Client {
	ctx := background
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func getGitHubUser(client *github.Client) *github.User {
	user, _, err := client.Users.Get(background, "")
	if err != nil {
		log.Println(err)
		return nil
	}
	return user
}