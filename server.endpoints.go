package main

import (
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/google/go-github/v56/github"
)

const port = "3639"

var indexPageData = IndexPageData{githubPublicID}

// Using Access Token and GitHub SDK can facilitate the use of GitHub API directly to structs.
type BasicPageData struct {
	ClientID       string
	User           github.User
	Mutuals        []MetaFollow
	IDontFollow    []MetaFollow
	TheyDontFollow []MetaFollow
}

// IndexPageData is the data for the index page template
type IndexPageData struct {
	ClientID string
}

type UnknownError struct {
	Err string
}

func serveWebApp() {
	log.Println("Serving Web App on port: ", port)
	http.HandleFunc("/", Index)
	http.HandleFunc("/result", Result)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

// The Index function renders the index page template and sends it as a response to the client.
func Index(w http.ResponseWriter, r *http.Request) {
	indexPage := template.Must(template.New("index.html").ParseFiles("./views/index.html"))
	if err := indexPage.Execute(w, indexPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderErrorPage(w http.ResponseWriter, err error) {
	render := template.Must(template.New("error.html").ParseFiles("./views/error.html"))
	htmlErr := UnknownError{err.Error()}
	if err := render.Execute(w, htmlErr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("Error: ", err)
}

// The Result function renders the result page template and sends it as a response to the client.
func Result(w http.ResponseWriter, r *http.Request) {
	// I need to abstract my GitHUB OAuth2.0 API call to a function
	// using the token, I need to display the result.
	// But Using HTMX, I can display the result on the same page.
	// Need to figure out what happens when I make the make the callback to the same page.
	if !r.URL.Query().Has("code") {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	accessKeys := getAccessToken(w, r)
	// bug: when user tries to refresh the page which has the same session token.
	// But the session token can only be used once.
	// The github API should return an error if the token is used more than once but it return 200 for every request.
	// So in the 2nd refresh, the access token is empty and the user is redirected to the login page.
	if accessKeys.AccessToken == "" {
		http.Redirect(w, r, "https://github.com/login/oauth/authorize?scope=user:follow&read:user&client_id="+githubPublicID, http.StatusTemporaryRedirect)
		return
	}
	ghClient := getGitHubClient(&accessKeys.AccessToken)
	ghUser := getGitHubUser(ghClient)

	// Get the followers of the user
	followers, err := GETFollowers(ghClient, *ghUser)
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	// Get the following of the user
	following, err := GETFollowing(ghClient, *ghUser)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
	/* The increased capacity of channel avoids deadlock for the c variable.
	The 3 go routines can run in concurrently without blocking each other.
	*/
	// results, resultsCh := make([][]MetaFollow, 3), make([]chan []MetaFollow, 3)
	resultsCh := make([]chan []MetaFollow, 3)
	for i := range resultsCh {
		resultsCh[i] = make(chan []MetaFollow, 1)
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go Mutuals(followers, following, resultsCh[0], &wg)
	go IDontFollow(followers, following, resultsCh[1], &wg)
	go TheyDontFollow(followers, following, resultsCh[2], &wg)
	wg.Wait()
	close(resultsCh[0])
	close(resultsCh[1])
	close(resultsCh[2])
	basicPageData := BasicPageData{githubPublicID, *ghUser, <-resultsCh[0], <-resultsCh[1], <-resultsCh[2]}
	render := template.Must(template.New("basic.html").ParseFiles("./views/basic.html"))
	if err := render.Execute(w, basicPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Result Page Served for user: ", *ghUser.Login)
}
