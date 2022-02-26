package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)


func main() {
	if !checkDb() {
		err := initDb()
		if err != nil {
			return
		}
	}
	r := mux.NewRouter()
	r.Use(handlers.ProxyHeaders)
	addCommentRouter(r)
	addShortenerRouter(r)
	http.ListenAndServe(":" + strconv.Itoa(ListenPort), r)

}