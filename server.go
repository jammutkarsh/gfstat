package main

import (
	"html/template"
	"net/http"
)



func Success(w http.ResponseWriter, b *BasicPageData) {

	// I just need to populate this struct with the data from the GitHub API
	basicPageData := BasicPageData{}

	render := template.Must(template.New("basic.tmpl").ParseFiles("views/basic.tmpl"))
	if err := render.Execute(w, basicPageData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


