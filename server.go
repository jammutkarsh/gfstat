package main

import (
	"html/template"
	"net/http"
)



func Success(w http.ResponseWriter, r *http.Request) {

	basicPageData := BasicPageData{User: user}

	render := template.Must(template.New("basic.tmpl").ParseFiles("views/basic.tmpl"))
	if err := render.Execute(w, basicPageData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


