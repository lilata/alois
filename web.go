package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)
func webIndex(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/web.html"))
	t.Execute(w, nil)
}
func webShortener(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/shortener.html"))
	t.Execute(w, nil)
}
func addWebRouter(r *mux.Router) {
	webRouter := r.PathPrefix("/web").Subrouter()
	webRouter.HandleFunc("", webIndex)
	webRouter.HandleFunc("/shortener", webShortener)

}