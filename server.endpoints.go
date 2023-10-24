package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/google/go-github/v56/github"
)

const port = "3639"

var indexPageData = IndexPageData{githubPublicID}

func init() {
	// check for env vars
	if githubPublicID == "" || githubServerSecret == "" {
		panic("GitHub OAuth2.0 Client ID and Secret ID not set")
	}
}

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
	indexPage := template.Must(template.New("index.html").ParseFiles("views/index.html"))
	if err := indexPage.Execute(w, indexPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("Index Page Served")
}

// The Result function renders the result page template and sends it as a response to the client.
func Result(w http.ResponseWriter, r *http.Request) {
	// I need to abstract my GitHUB OAuth2.0 API call to a function
	// using the token, I need to display the result.
	// But Using HTMX, I can display the result on the same page.
	// Need to figure out what happens when I make the make the callback to the same page.
	if !r.URL.Query().Has("code")  {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	accessKeys := getAccessToken(w, r)
	client := getGitHubClient(&accessKeys.AccessToken)
	user := getGitHubUser(client)

	// Get the followers of the user
	followers, err := GETFollowers(client, *user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the following of the user
	following, err := GETFollowing(client, *user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the mutuals of the user
	mutuals := Mutuals(followers, following)
	iDontFollow := FollowersYouDontFollow(followers, following)
	theyDontFollow := FollowingYouDontFollow(followers, following)
	basicPageData := BasicPageData{*user, mutuals, iDontFollow, theyDontFollow}
	render := template.Must(template.New("basic.html").ParseFiles("views/basic.html"))
	if err := render.Execute(w, basicPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Result Page Served")
}
