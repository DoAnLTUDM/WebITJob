package main

import (
    "database/sql"
    "log"
    "github.com/lib/pq"
    "strings"
    "reflect"
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

func getIdCompanyByJob(idJob int) (id int, err error) {
    rows, err:= DB.Query("select idComp from job where id=$1;", idJob)
    if err != nil {
        return id, err
    }
    if rows.Next() {
        err = rows.Scan(&id)
    }
    return id, nil

}

func getCompanyByName(name string) (company Company, err error)  {
    rows, err := DB.Query("select id, namecomp, address, country, logo, banner, intro, idSkill from company" +
        " where namecomp = $1;", name)
    if err != nil {
        return company, err
    }
    test := make(map[string]interface{})
    if rows.Next(){
        var idSkill_int64 []int64
        company = Company{}
        err = rows.Scan(&company.Id, &company.Name, &company.Address, &company.Country, &company.Logo, &company.Banner,
            (*pq.StringArray)(&company.Introduce), (*pq.Int64Array)(&idSkill_int64))
        for j:=0; j<len(idSkill_int64);j++{
            nameSkill, _ := getSkillName(idSkill_int64[j])
            company.Skill_name = append(company.Skill_name, nameSkill)
        }
        test["title"] = company.Introduce[0]
        test["content"] = company.Introduce[1:]
        company.Introduce_map = test
    }
    return company, err
}

func getLimitCompany(limit int) (companies []Company, err error){
    rows, err := DB.Query("select id, nameComp, address, country, logo, banner from company limit $1", limit)
    if err != nil {
        return companies, err
    }
    for rows.Next() {
        company := Company{}
        err = rows.Scan(&company.Id, &company.Name, &company.Address, &company.Country, &company.Logo, &company.Banner)
        if err != nil {
            return companies, err
        }
        company.Numjobs, _ = countJobs(company.Id)
        companies = append(companies, company)
    }
    return companies, nil
}

func getLimitJob(limit int) (jobs []Job, err error){
    rows, err := DB.Query("select job.id, job.title, job.salary, job.address, job.time_posted, job.description, job.skill," +
        " company.logo, job.idSkill from job inner join company on job.idComp = company.id limit $1;", limit)
    if err != nil {
        return jobs, err
    }

    for rows.Next() {
        job := Job{}
        var idSkill_int64 []int64
        err = rows.Scan(&job.Id, &job.Title, &job.Salary, &job.Address, &job.Time_posted, (*pq.StringArray)(&job.Job_description),
                (*pq.StringArray)(&job.Skill_expirence), &job.ImgComp, (*pq.Int64Array)(&idSkill_int64))
        var temp []string
        for i:=1; i<3;i++ {
            if (job.Job_description[i] != ""){
                temp = append(temp, job.Job_description[i])
            }
        }
        job.Job_description = temp
        for j:=0; j<len(idSkill_int64);j++{
            nameSkill, _ := getSkillName(idSkill_int64[j])
            job.Skill_name = append(job.Skill_name, nameSkill)
        }
        if err != nil {
            return jobs, err
        }
        jobs = append(jobs, job)
    }
    return jobs, nil
}

func countJobs(id_company int)  (id int, err error){
    rows, err := DB.Query("select count(idcomp) from job where idcomp=$1;",id_company)

    if err != nil {
        return id, err
    }
    if rows.Next() {
        err = rows.Scan(&id)
    }
    return id, err
}

func getSkillName(id int64) (name string, err error){
    rows, err := DB.Query("select nameskill from skill where id=$1;", int(id))
    if err != nil {
        return name, err
    }
    if rows.Next() {
        err = rows.Scan(&name)
    }
    return name, err
}

func getJobBySkill(val interface{}) (jobs []Job, err error){
    var id int
    var rows *sql.Rows
    if (reflect.TypeOf(val).Kind() == reflect.String ){
        rows, err := DB.Query("select id from skill where lower(nameskill)=$1;", strings.ToLower(val.(string)))
        if err != nil {
            return jobs, err
        }
        if rows.Next() {
            err = rows.Scan(&id)
        }
    }else{
        id = val.(int)
    }

    rows, err = DB.Query("select id, title, salary, address, time_posted, reason, description, skill, idcomp, idSkill " +
        "from job where $1 = ANY(idSkill);", id)
    if err != nil {
        return jobs, err
    }
    for rows.Next() {
        job := Job{}
        var idSkill_int64 []int64
        err = rows.Scan(&job.Id, &job.Title, &job.Salary, &job.Address, &job.Time_posted,
            (*pq.StringArray)(&job.Job_reason), (*pq.StringArray)(&job.Job_description),
            (*pq.StringArray)(&job.Skill_expirence), &job.IdComp, (*pq.Int64Array)(&idSkill_int64))
        if err != nil {
            return jobs, err
        }
        var temp []string
        for i:=1; i<3;i++ {
            if (job.Job_description[i] != ""){
                temp = append(temp, job.Job_description[i])
            }
        }
        job.Job_description = temp
        for j:=0; j<len(idSkill_int64);j++{
            nameSkill, _ := getSkillName(idSkill_int64[j])
            job.Skill_name = append(job.Skill_name, nameSkill)
        }
        img, _ := DB.Query("select logo from company where id=$1;", job.IdComp)
        if img.Next(){
            err = img.Scan(&job.ImgComp)
        }
        jobs = append(jobs, job)
    }
    return jobs, nil
}

func getJobByCompany(idcomp int) (jobs []Job, err error){
    rows, err := DB.Query("select id, title, salary, address, time_posted, reason, description, skill, idcomp, idSkill " +
        "from job where idcomp = $1;", idcomp)
    if err != nil {
        return jobs, err
    }
    for rows.Next() {
        job := Job{}
        var idSkill_int64 []int64
        err = rows.Scan(&job.Id, &job.Title, &job.Salary, &job.Address, &job.Time_posted,
            (*pq.StringArray)(&job.Job_reason), (*pq.StringArray)(&job.Job_description),
            (*pq.StringArray)(&job.Skill_expirence), &job.IdComp, (*pq.Int64Array)(&idSkill_int64))
        if err != nil {
            return jobs, err
        }
        var temp []string
        for i:=1; i<3;i++ {
            if (job.Job_description[i] != ""){
                temp = append(temp, job.Job_description[i])
            }
        }
        job.Job_description = temp
        for j:=0; j<len(idSkill_int64);j++{
            nameSkill, _ := getSkillName(idSkill_int64[j])
            job.Skill_name = append(job.Skill_name, nameSkill)
        }
        img, _ := DB.Query("select logo from company where id=$1;", job.IdComp)
        if img.Next(){
            err = img.Scan(&job.ImgComp)
        }
        jobs = append(jobs, job)
    }
    return jobs, nil
}

func getJobDetail(idjob int) (job Job, err error){
    DB, err = sql.Open("postgres", "user=postgres dbname=gwp password=a sslmode=disable")

    rows, err := DB.Query("select job.title, job.salary, job.address, job.time_posted, job.reason, job.description" +
        ", job.skill, company.logo, job.idSkill from job inner join company on job.idComp = company.id where job.id = $1;",
            idjob)
    if err != nil {
        return job, err
    }
    for rows.Next() {
        var idSkill_int64 []int64
        err = rows.Scan(&job.Title, &job.Salary, &job.Address, &job.Time_posted,
            (*pq.StringArray)(&job.Job_reason), (*pq.StringArray)(&job.Job_description),
            (*pq.StringArray)(&job.Skill_expirence), &job.ImgComp, (*pq.Int64Array)(&idSkill_int64))
        if err != nil {
            return job, err
        }
        for j:=0; j<len(idSkill_int64);j++ {
            nameSkill, _ := getSkillName(idSkill_int64[j])
            job.Skill_name = append(job.Skill_name, nameSkill)
        }
    }

    return job, nil
}