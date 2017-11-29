package main

import (
    "html/template"
    "net/http"
    "log"
    "fmt"
)

func templateHomePage(w http.ResponseWriter, r *http.Request) {
    name := "main_page"
    t, err := template.ParseFiles("template/index.html", "template/main_page.html")
    if err != nil{
        log.Fatal(err)
    }
    t.Execute(w, name)
}

func notFoundPage(w http.ResponseWriter, r *http.Request){
    t, _ := template.ParseFiles("template/index.html", "template/404.html")
    t.Execute(w,"" )
}

func login(w http.ResponseWriter, r *http.Request){

}

func getRequestBody(w http.ResponseWriter, r *http.Request){
    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)
    fmt.Println(w, string(body))
}
