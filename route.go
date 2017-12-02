package main

import (
    "html/template"
    "net/http"
    "log"
    "fmt"
)

func templateHomePage(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("template/index.html", "template/main_page.html",
        "template/job_recommend.html", "template/company_recommend.html")
    if err != nil{
        log.Fatal(err)
    }
    company, err := getLimitCompany(9)
    if err != nil{
        log.Fatal(err)
    }
    log.Println(company)
    t.Execute(w, company)
}

func notFoundPage(w http.ResponseWriter, r *http.Request){
    t, _ := template.ParseFiles("template/index.html", "template/404.html")
    t.Execute(w,nil)
}

func login(w http.ResponseWriter, r *http.Request){

}

func getRequestBody(w http.ResponseWriter, r *http.Request){
    len := r.ContentLength
    body := make([]byte, len)
    r.Body.Read(body)
    fmt.Println(w, string(body))
}
