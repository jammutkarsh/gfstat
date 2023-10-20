package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const port = "3639"

var (
	// GitHub OAuth Config
	clientId     = os.Getenv("GH_BASIC_CLIENT_ID") // like public key
	clientSecret = os.Getenv("GH_BASIC_SECRET_ID") // like private key
	// Frontend
	// basicPage     = template.Must(template.New("basic.tmpl").ParseFiles("views/basic.tmpl"))
	indexPageData = IndexPageData{clientId}
	// Context
	background = context.Background()
)

type IndexPageData struct {
	ClientId string
}

type BasicPageData struct {
	User   *github.User
}

type Access struct {
	AccessToken string `json:"access_token"`
	Scope       string
}

func serve() {
	fmt.Println("http://127.0.0.1:3639")
	http.HandleFunc("/", Index)
	// http.HandleFunc("/success", Success)
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

// func Success(w http.ResponseWriter, r *http.Request) {
// 	accessToken := LoginWithGitHub(r)
// 	client := getGitHubClient(accessToken)

// 	user, _, err := client.Users.Get(background, "")
// 	if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	basicPageData := BasicPageData{user}

// 	render := template.Must(template.New("basic.html").ParseFiles("views/basic.html"))
// 	if err := render.Execute(w, basicPageData); err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func LoginWithGitHub(w http.ResponseWriter, r *http.Request) (accessToken string) {
// 	code := r.URL.Query().Get("code")
// 	values := url.Values{"client_id": {clientId}, "client_secret": {clientSecret}, "code": {code}, "accept": {"json"}}

// 	req := fasthttp.AcquireRequest()
// 	resp := fasthttp.AcquireResponse()
// 	defer fasthttp.ReleaseRequest(req)
// 	defer fasthttp.ReleaseResponse(resp)

// 	req.SetRequestURI("https://github.com/login/oauth/access_token")
// 	req.SetBodyString(values.Encode())
// 	req.Header.Set("Accept", "application/json")
// 	req.Header.SetMethod("POST")

// 	if err := fasthttp.Do(req, resp); err != nil {
// 		ctx.Error(err.Error(), fasthttp.StatusPreconditionFailed)
// 		return
// 	}
// 	if resp.StatusCode() != fasthttp.StatusOK {
// 		ctx.Error(fmt.Sprintf("Retrieving access token failed: %d", resp.StatusCode()), resp.StatusCode())
// 		return
// 	}
// 	var access Access

// 	if err := json.Unmarshal(resp.Body(), &access); err != nil {
// 		ctx.Error(fmt.Sprintf("JSON-Decode-Problem: %s", err), fasthttp.StatusInternalServerError)
// 		return
// 	}

// 	if access.Scope != "user:email" {
// 		ctx.Error(fmt.Sprintf("Wrong token scope: %s", access.Scope), fasthttp.StatusPreconditionFailed)
// 		return
// 	}

// 	return access.AccessToken
// }

// Authenticates GitHub Client with provided OAuth access token
func getGitHubClient(accessToken string) *github.Client {
	ctx := background
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
