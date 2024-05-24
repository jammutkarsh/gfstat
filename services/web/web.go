package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/JammUtkarsh/gfstat/services/client"
	"github.com/JammUtkarsh/gfstat/services/core"
	"github.com/google/go-github/github"
)

type ResultData struct {
	ClientID       string
	User           github.User
	Mutuals        []core.MetaFollow
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
	// The github API should return an error if the token is used more than once but it return 200 for every request.
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
		err       error
	)

	followers, err = client.GETFollowers(ghClient, *ghUser, *ghResp)
	if err != nil {
		renderErrorPage(w, err)
		return
	}
	following, err = client.GETFollowing(ghClient, *ghUser, *ghResp)
	if err != nil {
		renderErrorPage(w, err)
		return
	}

	mutuals := core.Mutuals(followers, following)
	iDontFollow := core.IDontFollow(followers, following)
	theyDontFollow := core.TheyDontFollow(followers, following)

	basicPageData := ResultData{client.GithubPublicID, *ghUser, mutuals, iDontFollow, theyDontFollow}
	render := template.Must(template.New("basic.html").ParseFiles("./views/basic.html"))
	if err = render.Execute(w, basicPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Result Page Served for user: ", *ghUser.Login)
}
