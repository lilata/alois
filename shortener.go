package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)
func genAlias() string {
	chars := []rune(
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			"0123456789" +
			"-._~",
	)
	rand.Seed(time.Now().UnixNano())
	var runes []rune
	for i := 0; i < URLAliasLength; i++ {
		runes = append(runes, chars[rand.Intn(len(chars))])
	}
	return string(runes)
}
func redirectToUrlAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias, ok := vars["alias"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	db := getDb()
	defer db.Close()
	row, err := db.Query("SELECT value FROM url_aliases WHERE alias = ? LIMIT 1", alias)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	var value string
	for row.Next() {
		row.Scan(&value)
	}
	if len(value) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err = url.ParseRequestURI(value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", value)
	w.WriteHeader(http.StatusFound)

}
func addUrlAlias(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	toShorten := r.Form.Get("url")
	if _, err := url.ParseRequestURI(toShorten); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db := getDb()
	defer db.Close()
	row, err := db.Query("SELECT alias FROM url_aliases WHERE value = ? LIMIT 1", toShorten)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var alias string
	for row.Next() {
		row.Scan(&alias)
	}
	scheme := getUrlScheme(r)
	host := r.Host
	if (len(ShortenerHost) != 0) {
		host = ShortenerHost
	}
	if len(alias) != 0 {
		w.Write([]byte(fmt.Sprintf("%s://%s/s/%s", scheme, host, alias)))
		return
	}
	cnt := 1
	for cnt > 0 {
		alias = genAlias()
		err := db.QueryRow("SELECT COUNT(*) FROM url_aliases WHERE alias = ? LIMIT 1", alias).Scan(&cnt)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}


	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stmt, err := tx.Prepare("INSERT INTO url_aliases(alias, value) VALUES (?, ?)")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(alias, toShorten)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tx.Commit()
	w.Write([]byte(fmt.Sprintf("%s://%s/s/%s", scheme, host, alias)))

}

func addShortenerRouter(r *mux.Router) {
	shortenerRouter := r.PathPrefix("/shortener").Subrouter()
	shortenerRouter.HandleFunc("/add", addUrlAlias)
	r.HandleFunc("/s/{alias}", redirectToUrlAlias)
}