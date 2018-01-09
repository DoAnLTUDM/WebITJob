package main

import (
    "html/template"
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "strconv"
)

func templateHomePage(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("template/index.html", "template/main_page.html", "template/recommend.html",
        "template/listjob.html")
    if err != nil{
        log.Println(err)
    }
    company, err := getLimitCompany(9)
    if err != nil{
        log.Println(err)
    }

    job, err := getLimitJob(3)
    if err != nil{
        log.Println(err)
    }
    map_new := map[string]interface{}{"jobs":job, "companies": company}

    t.Execute(w, map_new)
}

func notFoundPage(w http.ResponseWriter, r *http.Request){
    t, _ := template.ParseFiles("template/index.html", "template/404.html")
    t.Execute(w,nil)
}

func login(w http.ResponseWriter, r *http.Request){

}

func infoJobPage(w http.ResponseWriter, r *http.Request){
    val := mux.Vars(r)
    idstr := val["id"]
    idjob, _:= strconv.Atoi(idstr)
    idcomp, _:= getIdCompanyByJob(idjob)
	t, _ := template.ParseFiles("template/index.html", "template/infojob.html", "template/listjob.html")
    jobs, _ := getJobBySkill("java")
    jodDetail, _:= getJobDetail(idcomp)
    map_new := map[string]interface{}{"job-detail": jodDetail,"jobs":jobs}
    t.Execute(w,map_new)
}

func infoCompanyPage(w http.ResponseWriter, r *http.Request){
    val := mux.Vars(r)
    company_name := val["company_name"]
	t, _ := template.ParseFiles("template/index.html", "template/infocompany.html",  "template/listjob.html")
    company, _ := getCompanyByName(company_name)
    job, _ := getJobByCompany(company.Id)
    map_new := map[string]interface{}{ "company": company, "jobs":job}
	t.Execute(w,map_new)
}

func searchJobBySkill(w http.ResponseWriter, r *http.Request){
    t,_:= template.ParseFiles("template/index.html", "template/main_page.html",
        "template/listjob.html")
    vars := mux.Vars(r)
    skill := vars["skill"]
    job, err := getJobBySkill(skill)
    if err != nil{
        log.Println(err)
    }
    t.Execute(w,job)
}

func search(w http.ResponseWriter, r *http.Request){

}
