package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
)


func main() {
	r := mux.NewRouter()
    schedule_crawl_data()
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css/"))))
	r.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img/"))))
	r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js/"))))
	r.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts/", http.FileServer(http.Dir("./public/fonts/"))))

	r.HandleFunc("/homepage", templateHomePage)
	r.HandleFunc("/notfound", notFoundPage)
	r.HandleFunc("/infojob", infoJobPage)
	r.HandleFunc("/company/{company_name}",infoCompanyPage)
	r.HandleFunc("/jobskill/{skill}", searchJobBySkill)
	r.HandleFunc("/job/{title}/{id}", infoJobPage)
	log.Fatal(http.ListenAndServe(":8000",r))
}
