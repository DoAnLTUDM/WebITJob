package main

import (
    "database/sql"
    "log"
    "github.com/lib/pq"
)

var DB *sql.DB

func init(){
    var err error
    //Connect to Database
    DB, err = sql.Open("postgres", "user=postgres dbname=gwp password=a sslmode=disable")
    if err != nil{
        panic(err)
    }
}

func (company *Company) Create() (id int, err error)  {
    statement := "insert into company (nameComp, address, country, logo, banner, intro, idSkill) " +
        "values ($1, $2, $3, $4, $5, $6, $7) returning id;"
    stmtm, err := DB.Prepare(statement)
    defer stmtm.Close()
    if err != nil{
        log.Println(err)
        return id, err
    }
    stmtm.QueryRow(company.Name, company.Address, company.Country, company.Logo, company.Banner,
        pq.Array(company.Introduce), pq.Array(company.IdSkills)).Scan(&id)
    return id,nil
}

func (job *Job) Create() (id int, err error) {
    statement := "insert into job (idSkill, idComps, title, salary, address, time_posted, reason, description, skill)" +
        " values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;"
    stmtm, err := DB.Prepare(statement)
    defer stmtm.Close()
    if err != nil{
        log.Println(err)
        return id, err
    }
    stmtm.QueryRow(pq.Array(job.IdSkills), job.IdComp, job.Title, job.Salary, job.Address,
        job.Time_posted, pq.Array(job.Job_reason), pq.Array(job.Job_description), pq.Array(job.Skill_expirence)).Scan(&id)
    return id,nil
}

func (skill *Skill) Create() (id int, err error){
    statement :=  "insert into skill (nameSkill) values ($1) returning id;"
    stmtm, err := DB.Prepare(statement)
    defer stmtm.Close()
    if err != nil{
        log.Println(err)
        return id, err
    }
    stmtm.QueryRow(skill.Name).Scan(&id)
    return id,nil
}

func getIdSkill(name string) (id int, err error) {
    statement := "select id from skill where nameSkill=$1;"
    stmtm, err := DB.Prepare(statement)
    defer stmtm.Close()
    if err != nil {
        return id, err
    }
    stmtm.QueryRow(name).Scan(&id)
    return id, nil

}

func getIdCompany(name string) (id int, err error) {
    statement := "select id from company where nameComp=$1;"
    stmtm, err := DB.Prepare(statement)
    defer stmtm.Close()
    if err != nil {
        return id, err
    }
    stmtm.QueryRow(name).Scan(&id)
    return id, nil

}

//func getCompanyById(id int) (company Company, err error)  {
//    company = Company{}
//    err = DB.QueryRow("select company_name, address, country, logo, urlC from company where id = $1", id).Scan(
//        &company.Name, &company.Address, &company.Country, &company.Logo, &company.UrlC)
//    return company, err
//}

func getLimitCompany(limit int) (companies []Company, err error){
    rows, err := DB.Query("select id, nameComp, address, country, logo, banner, intro from company limit $1", limit)
    if err != nil {
        return companies, err
    }
    for rows.Next() {
        company := Company{}
        err = rows.Scan(&company.Id, &company.Name, &company.Address, &company.Country, &company.Logo, &company.Banner,
            (*pq.StringArray)(&company.Introduce))
        if err != nil {
            return companies, err
        }
        company.Numjobs, _ = countJobs(company.Id)
        companies = append(companies, company)
    }
    return companies, nil
}

func getLimitJob(limit int) (jobs []Job, err error){
    rows, err := DB.Query("select title, salary, address, time_posted, reason," +
        "description, skill from job limit $1", limit)
    if err != nil {
        return jobs, err
    }
    for rows.Next() {
        job := Job{}
        err = rows.Scan(&job.Title, &job.Salary, &job.Address, &job.Time_posted,
            (*pq.StringArray)(&job.Job_reason), (*pq.StringArray)(&job.Job_description),
                (*pq.StringArray)(&job.Skill_expirence))
        if err != nil {
            return jobs, err
        }
        jobs = append(jobs, job)
    }
    return jobs, nil
}

func getJobBySkill(name string) (jobs []Job, err error){
    var id int
    rows, err := DB.Query("select id from skill where lower(nameskill)=$1;", name)
    if err != nil {
        return jobs, err
    }
    if rows.Next() {
        err = rows.Scan(&id)
    }

    rows, err = DB.Query("select title, salary, address, time_posted, reason, description, skill, idcomp " +
        "from job where $1 = ANY(idSkill);", id)
    if err != nil {
        return jobs, err
    }
    for rows.Next() {
        job := Job{}
        err = rows.Scan(&job.Title, &job.Salary, &job.Address, &job.Time_posted,
            (*pq.StringArray)(&job.Job_reason), (*pq.StringArray)(&job.Job_description),
            (*pq.StringArray)(&job.Skill_expirence), &job.IdComp)
        if err != nil {
            return jobs, err
        }
        img, _ := DB.Query("select logo from company where id=$1;", job.IdComp)
        if img.Next(){
            err = img.Scan(&job.ImgComp)
        }
        jobs = append(jobs, job)
    }
    return jobs, nil
}

func countJobs(id_company int)  (id int, err error){
    log.Print(id_company)
    rows, err := DB.Query("select count(idcomp) from job where idcomp=$1;",id_company)

    if err != nil {
        return id, err
    }
    if rows.Next() {
        err = rows.Scan(&id)
    }
    log.Print(id)
    return id, err
}