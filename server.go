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

var (
	// GitHub OAuth Config
	clientId     = os.Getenv("GH_BASIC_CLIENT_ID") // like public key
	clientSecret = os.Getenv("GH_BASIC_SECRET_ID") // like private key
	// Frontend
	indexPageData = IndexPageData{clientId}
	// Context
	background = context.Background()
)

// To get Secret ID
type IndexPageData struct {
	ClientId string
}

// Using Access Token and GitHub SDK can facilitate the use of GitHub API directly to structs.
type BasicPageData struct {
	User   *github.User
	Mutuals []MetaFollow
	iDontFollow []MetaFollow
	theyDontFollow []MetaFollow
}

type Access struct {
	AccessToken string `json:"access_token"`
	Scope       string
}

func serve() {
	fmt.Println("http://127.0.0.1:3639")
	http.HandleFunc("/", Index)
	http.HandleFunc("/success", Success)
	if err:=http.ListenAndServe(":"+port, nil); err!=nil {
		panic(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	indexPage := template.Must(template.New("index.tmpl").ParseFiles("views/index.tmpl"))
	if err := indexPage.Execute(w, indexPageData); err != nil {
		log.Println(err)
	}
}

func Success(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	values := url.Values{"client_id": {clientId}, "client_secret": {clientSecret}, "code": {code}, "accept": {"json"}}

	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(values.Encode()))
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	var access Access

	if err := json.NewDecoder(resp.Body).Decode(&access); err != nil {
		log.Println("JSON-Decode-Problem: ", err)
		return
	}

	if access.Scope != "user:email" {
		log.Println("Wrong token scope: ", access.Scope)
		return
	}

	client := getGitHubClient(access.AccessToken)

	user, _, err := client.Users.Get(background, "")
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	basicPageData := BasicPageData{User: user}

	render := template.Must(template.New("basic.tmpl").ParseFiles("views/basic.tmpl"))
	if err := render.Execute(w, basicPageData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
