package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

func main() {
	schedule_crawl_data()
	r := mux.NewRouter()
	cssHandler := http.FileServer(http.Dir("./public/css/"))
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	imagesHandler := http.FileServer(http.Dir("./public/img/"))
	http.Handle("/img/", http.StripPrefix("/img/", imagesHandler))
	jsHandler := http.FileServer(http.Dir("./public/js/"))
	http.Handle("/js/", http.StripPrefix("/js/", jsHandler))


	r.HandleFunc("/homepage", templateHomePage)
	r.HandleFunc("/notfound", notFoundPage)
	log.Println()
	http.Handle("/homepage",r)
	http.Handle("/notfound",r)
	log.Fatal(http.ListenAndServe(":8000",nil))
}

