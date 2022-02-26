package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"
)
type Comment struct {
	Username string
	Content string
	Time    int64
}
func (t *Comment) StringTime() string {
	return time.Unix(t.Time, 0).Format(time.UnixDate)
}
func addComment(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	username := r.Form.Get("username")
	content := r.Form.Get("content")
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	redirectTo := fmt.Sprintf("%s://%s/comment/render?key=%s", getUrlScheme(r), r.Host, key)
	if len(username) == 0 || len(content) == 0 {
		redirectTo = fmt.Sprintf("%s&submit_error=1", redirectTo)
		w.Header().Set("Location", redirectTo)
		w.WriteHeader(http.StatusFound)
		return
	}
	t := time.Now().Unix()
	db := getDb()
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}
	stmt, err := tx.Prepare("INSERT INTO comments(key, username, content, time) " +
		"VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(key, username, content, t)
	if err != nil {
		log.Println(err)
	}
	tx.Commit()
	w.Header().Set("Location", redirectTo)
	w.WriteHeader(http.StatusFound)
}
func renderComments(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	submitError := r.Form.Get("submit_error")
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db := getDb()
	defer db.Close()
	rows, err := db.Query("SELECT username, content, time FROM comments WHERE key=? ORDER BY time DESC", key)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var c Comment
		_ = rows.Scan(&c.Username, &c.Content, &c.Time)
		comments = append(comments, c)
	}
	t := template.Must(template.ParseFiles("templates/comments.html"))
	_, err = url.ParseRequestURI(key)
	keyIsUrl := err == nil
	commentUrl := fmt.Sprintf("%s://%s/comment/add", getUrlScheme(r), r.Host)
	data := struct {
		Key string
		Comments []Comment
		KeyIsUrl bool
		CommentUrl string
		SubmitError bool
	} {
		Key: key,
		Comments: comments,
		KeyIsUrl: keyIsUrl,
		CommentUrl: commentUrl,
		SubmitError: submitError == "1",
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
	}

}

func addCommentRouter(r *mux.Router) {
	commentRouter := r.PathPrefix("/comment").Subrouter()
	commentRouter.HandleFunc("/add", addComment)
	commentRouter.HandleFunc("/render", renderComments)
}