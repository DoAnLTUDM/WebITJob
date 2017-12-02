package main

import (
    "github.com/gorilla/mux"
    "net/http"
    "log"
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
    fontsHandler := http.FileServer(http.Dir("./public/fonts/"))
    http.Handle("/fonts/", http.StripPrefix("/fonts/", fontsHandler))

    r.HandleFunc("/homepage", templateHomePage)
    r.HandleFunc("/notfound", notFoundPage)
    r.HandleFunc("/infojob", infoJobPage)
    r.HandleFunc("/infocompany",infoCompanyPage)
    log.Println()
    http.Handle("/homepage",r)
    http.Handle("/notfound",r)
    http.Handle("/infojob",r)
    http.Handle("/infocompany",r)
    log.Fatal(http.ListenAndServe(":8000",nil))
}
