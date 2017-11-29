package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"log"
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

func (company *Company) Create() (err error)  {
	statement := "insert into company (company_name, address, country, logo, urlC) " +
		"values ($1, $2, $3, $4, $5) returning id;"
	stmtm, err := DB.Prepare(statement)
	if err != nil{
		log.Println(err)
		return err
	}
	defer stmtm.Close()
	stmtm.QueryRow(company.Name, company.Address, company.Country, company.Logo, company.UrlC)
	return nil
}

func (job *Job) Create() (err error) {
	statement := "insert into job (title, salary, address, time_posted, reason, description," +
		" skill, qualification, company_name) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id;"
	stmtm, err := DB.Prepare(statement)
	if err != nil{
		log.Println(err)
		return err
	}
	defer stmtm.Close()
	stmtm.QueryRow(job.Title, job.Salary, job.Address, job.Time_posted, job.Job_reason, job.Job_description,
		job.Skill_expirence, job.Qualification, job.Company_name )
	return nil
}

func getCompanyById(id int) (company Company, err error)  {
	company = Company{}
	err = DB.QueryRow("select company_name, address, country, logo, urlC from company where id = $1", id).Scan(
		&company.Name, &company.Address, &company.Country, &company.Logo, &company.UrlC)
	return company, err
}

func getLimitCompany(limit int) (companies []Company, err error){
	rows, err := DB.Query("select company_name, urlC, logo, address, country from company limit $1", limit)
	if err != nil {
		return companies, err
	}
	for rows.Next() {
		company := Company{}
		err = rows.Scan(&company.Name, &company.UrlC, &company.Logo, &company.Address, &company.Country)
		if err != nil {
			return companies, err
		}
		companies = append(companies, company)
	}
	return companies, nil
}