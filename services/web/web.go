package web

import (
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/JammUtkarsh/gfstat/services/client"
	"github.com/JammUtkarsh/gfstat/services/core"
	"github.com/google/go-github/github"
)

// ResultData WARN: Changing name of the variables here directly affects HTML generation.
// The changes made here should be synchronized with view/*.html files
type ResultData struct {
	ClientID       string
	User           github.User
	Mutual         []core.MetaFollow
	IDontFollow    []core.MetaFollow
	TheyDontFollow []core.MetaFollow
}

// IndexData is the data for the index page template
type IndexData struct {
	ClientID string
}

const port = "3639"

var indexPageData = IndexData{client.GithubPublicID}

func ServeWebApp() {
	log.Println("Serving Web App on port: ", port)
	http.HandleFunc("/", index)
	http.HandleFunc("/result", result)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	indexPage := template.Must(template.New("index.html").ParseFiles("./views/index.html"))
	if err := indexPage.Execute(w, indexPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderErrorPage(w http.ResponseWriter, err error) {
	render := template.Must(template.New("error.html").ParseFiles("./views/error.html"))
	if err := render.Execute(w, err.Error()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("Error: ", err)
}

func result(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("code") {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	accessKeys := client.GetAccessToken(w, r)
	// bug: when user tries to refresh the page which has the same session token.
	// But the session token can only be used once.
	// The GitHub API should return an error if the token is used more than once, but it returns 200 for every request.
	// So in the 2nd refresh, the access token is empty and the user is redirected to the login page.
	if accessKeys.AccessToken == "" {
		http.Redirect(w, r, "https://github.com/login/oauth/authorize?scope=user:follow&read:user&client_id="+client.GithubPublicID, http.StatusTemporaryRedirect)
		return
	}
	ghClient := client.GetGitHubClient(&accessKeys.AccessToken)
	ghUser, ghResp := client.GetGitHubUser(ghClient)

	// Get the followers of the user
	var (
		followers []core.MetaFollow
		following []core.MetaFollow
		wg        sync.WaitGroup
		err       error
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		followers, err = client.GETFollowers(ghClient, *ghUser, *ghResp)
	}()
	go func() {
		defer wg.Done()
		following, err = client.GETFollowing(ghClient, *ghUser, *ghResp)
	}()
	wg.Wait()

	if err != nil {
		renderErrorPage(w, err)
		return
	}

	resultsCh := make([]chan []core.MetaFollow, 3)
	for i := range resultsCh {
		resultsCh[i] = make(chan []core.MetaFollow, 1)
	}
	wg.Add(3)
	go core.Mutuals(followers, following, resultsCh[0], &wg)
	go core.IDontFollow(followers, following, resultsCh[1], &wg)
	go core.TheyDontFollow(followers, following, resultsCh[2], &wg)
	wg.Wait()

	basicPageData := ResultData{client.GithubPublicID, *ghUser, <-resultsCh[0], <-resultsCh[1], <-resultsCh[2]}
	render := template.Must(template.New("basic.html").ParseFiles("./views/basic.html"))
	if err = render.Execute(w, basicPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Result Page Served for user: ", *ghUser.Login)
}
