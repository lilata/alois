package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
)
func pipedURIs() []string {
	resp, err := http.Get("https://piped-instances.kavin.rocks/")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var names []struct {
		Name string
	}
	_ = json.Unmarshal(body, &names)
	var uris []string
	for _, name := range names {
		uri := strings.TrimSuffix(name.Name, " (Official)")
		uris = append(uris, uri)
	}
	return uris
}
func invidiousURIs() []string{
	resp, err := http.Get("https://api.invidious.io/instances.json")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var uris [][]string
	_ = json.Unmarshal(body, &uris)
	var clearWebUris []string
	for _, uri := range uris {
		if !strings.HasSuffix(uri[0], "onion") {
			clearWebUris = append(clearWebUris, "https://" + uri[0])
		}
	}
	return clearWebUris
}
func watchYoutube(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	watchId := r.Form.Get("v")
	uris := append(invidiousURIs(), pipedURIs()...)
	idx := rand.Intn(len(uris))
	w.Header().Set("Location", fmt.Sprintf("%s/watch?v=%s", uris[idx], watchId))
	w.WriteHeader(http.StatusFound)
}
func addYoutubeRouter(r *mux.Router) {
	ytbRouter := r.PathPrefix("/ytb").Subrouter()
	ytbRouter.HandleFunc("/watch", watchYoutube)
}