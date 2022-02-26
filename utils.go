package main

import "net/http"

func getUrlScheme(r *http.Request) string {
	s := r.URL.Scheme
	if len(s) == 0 {
		s = "http"
	}
	return s
}