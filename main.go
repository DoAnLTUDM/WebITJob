package main

import (
	"html/template"
	"net/http"
)

func templateHomePage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, "")
}
func main() {
	server := http.Server{
		Addr: "120.0.0.1:8000",
	}
	http.HandleFunc("/homepage", templateHomePage)
	server.ListenAndServe()
}
