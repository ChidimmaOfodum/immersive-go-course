package main

import (
	"fmt"
	"html"
	"os"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			params := r.URL.Query()
			var list string
			for key, value := range params {
				stringifiedValue := fmt.Sprintf("%s", value)
				list += fmt.Sprintf("<li>%s: %s</li>\n", key, html.EscapeString(stringifiedValue))
			}
			w.Write([]byte(fmt.Sprintf("<!DOCTYPE html>\n<html>\n<em>Hello, world</em>\n<p>Query parameters: </p>\n<ul>%s</ul>", list)))
		}
	})
	http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok || (user != os.Getenv("AUTH_USERNAME") && password != os.Getenv("AUTH_PASSWORD")) {
			w.Header().Set("WWW-Authenticate", `Basic realm="localhost", charset="UTF-8"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(fmt.Sprintf("<!DOCTYPE html><html><p>Hello, %s\n </p>", user)))
	})

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK\n"))
	})

	http.Handle("/404", http.NotFoundHandler())

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
	})
	http.ListenAndServe(":8080", nil)

}
