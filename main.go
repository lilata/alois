package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
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
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))
	log.Println("Alois API starts...")
	http.ListenAndServe(":" + strconv.Itoa(ListenPort), r)

}