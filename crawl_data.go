package main

import (
	"log"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strings"
	"regexp"
	"golang.org/x/sync/errgroup"
	"encoding/json"
	"io/ioutil"
)

type Job struct {
	title string `json:"title"`
	jobId string `json:"jobId"`
	salary string `json:"salary"`
	address string `json:"address"`
	time_posted string `json:"time_posted"`
	job_reason string `json:"job_reason"`
	job_description string `json:"job_description"`
	skill_expirence string `json:"skill_expirence"`
	qualification string `json:"qualification"`
	company_name string `json:"company_name"`
}

type Jobs struct{
	totalJobs int `json:total_jobs`
	ListJobs []Job `json:"jobs"`
}

type Company struct{
	name string `json:"name"`
	address string `json:"address"`
	country string `json:"country"`
	logo string `json:"image"`
	urlC string `json:"url"`
}

type Companies struct{
	totalCompanies int `json:total_companies`
	ListCompanies []Company `json:"companies"`
}

func NewJobs() *Jobs  {
	return &Jobs{}
}

func NewCompanies() *Companies {
	return &Companies{}
}

func checkError(err error){
	if err != nil{
		log.Println(err)
	}
}

func crawl_data(url string)  {
	//Get Url
	if len(url) < 1{
		log.Print("Please specify start page")
		os.Exit(1)
	}
	currentUrl := url
	log.Println("Access to ",currentUrl)

	companies := NewCompanies()
	err := companies.getCompaniesByUrl(currentUrl)
	checkError(err)
	companiesJson, err := json.Marshal(companies) // Chuyển kiểu dữ liệu companies sang JSON
	checkError(err)
	err = ioutil.WriteFile("companies.json", companiesJson, 0644) // Ghi dữ liệu vào file JSON
	checkError(err)

	jobs := NewJobs()
	err = jobs.GetAllJobs(companies)
	checkError(err)
	jobsJson, err := json.Marshal(jobs) // Chuyển kiểu dữ liệu jobs sang JSON
	checkError(err)
	err = ioutil.WriteFile("jobs.json", jobsJson, 0644) // Ghi dữ liệu vào file JSON
	checkError(err)
}

func (companies *Companies) getCompaniesByUrl (urlC string) error{
	doc, err := goquery.NewDocument(urlC)
	if err != nil {
		return err
	}

	doc.Find("#container .top-companies .col-md-4 ").Each(func(i int, selection *goquery.Selection) {
		companyUrl,exist := selection.Find("a").Attr("href")
		if exist {
			if strings.Contains(urlC, "itviec"){
				companyUrl = "https://itviec.com" + companyUrl
			}
			if strings.Contains(urlC, "vietnamworks"){
				companyUrl = "https://www.vietnamworks.com" + companyUrl
			}
		} else {
			companyUrl = "#"
		}
		companies.getInformationCompanies(companyUrl)
	})

	return nil
}

func (companies *Companies) getInformationCompanies(companyUrl string) error  {
	doc, err := goquery.NewDocument(companyUrl)

	if err != nil {
		return err
	}

	doc.Find("#container .company-content").Each(func(i int, selection *goquery.Selection) {
		companyImg,_ := selection.Find(".headers .logo-container img").Attr("src")
		companyName := selection.Find(".headers .name-and-info .title").Text()
		companyAddr := selection.Find(".col-md-3 .map-address").Text()
		companyCountry := selection.Find(".headers .company-info .country span").Text()

		company := Company{
			name:companyName,
			urlC:companyUrl,
			logo:companyImg,
			address:companyAddr,
			country:companyCountry,
		}

		companies.totalCompanies++
		companies.ListCompanies = append(companies.ListCompanies, company)
	})
	return nil
}

func (jobs *Jobs) getJobsByUrl (urlJ string) error{
	doc, err := goquery.NewDocument(urlJ)
	if err != nil {
		return err
	}

	doc.Find(".job_content .title").Each(func(i int, selection *goquery.Selection) {
		jobUrl,exist := selection.Find("a").Attr("href")
		if exist{
			if strings.Contains(urlJ, "itviec"){
				jobUrl = "https://itviec.com" + jobUrl
			}
			if strings.Contains(urlJ, "vietnamworks"){
				jobUrl = "https://www.vietnamworks.com" + jobUrl
			}
		}else {
			jobUrl = "#"
		}
		jobs.getDetailJob(jobUrl)
	})

	return nil
}

func (jobs *Jobs) getDetailJob(jobUrl string) error  {
	doc, err := goquery.NewDocument(jobUrl)
	if err != nil {
		return err
	}
	doc.Find("#container").Each(func(i int, selection *goquery.Selection) {
		jobId,_ := selection.Find(".show_content").Attr("id")
		jobTitle := strings.TrimSpace(selection.Find(".job_title").Text())
		jobSalary := selection.Find(".salary_text").Text()
		jobCompany := selection.Find(".inside .employer-info a").Text()
		re := regexp.MustCompile(`\r?\n`)
		jobAddress := re.ReplaceAllString(strings.TrimSpace(
			selection.Find(".address__full-address").Text()), " ")
		jobTimePosted := strings.TrimSpace(selection.Find(".distance-time-job-posted").Text())
		var jobReason string
		s := selection.Find(".love_working_here .culture_description")
		s = s.Contents().Each(func(i int, selection *goquery.Selection) {

			//if goquery.NodeName(selection) == "li" || goquery.NodeName(selection) == "ul" {
			//	jobReason += selection.Find(goquery.NodeName(selection)).Text() + " *** "
			//	fmt.Println("++++",goquery.NodeName(selection),"++++")
			//}
			if   goquery.NodeName(selection) == "p" || goquery.NodeName(selection) == "div" {
				jobReason += selection.Text() + " *** "
			}
			s1 := selection.Find("li")
			s1 = s1.Contents().Each(func(i int, selection *goquery.Selection) {
				if   goquery.NodeName(selection) == "strong" {
					jobReason += selection.Text()
				}
				if goquery.NodeName(selection) == "#text"{
					jobReason += selection.Text() + " *** "
				}
			})
		})
		var jobDescription string
		s = selection.Find(".job_description .description li")
		s = s.Contents().Each(func(i int, selection *goquery.Selection) {

			if goquery.NodeName(selection) == "#text" {
				jobDescription += selection.Text() + " *** "
			}
		})

		var jobSkill string
		s = selection.Find(".skills_experience .experience ul li")
		s = s.Contents().Each(func(i int, selection *goquery.Selection) {
			if goquery.NodeName(selection) == "#text" {
				jobSkill += selection.Text() + " *** "
			}
		})

		//jobCompany := selection.Find("")
		//
		job := Job {
			jobId:jobId,
			title:jobTitle,
			company_name: jobCompany,
			salary: jobSalary,
			address: jobAddress,
			time_posted: jobTimePosted,
			job_reason: jobReason ,
			job_description: jobDescription,
			skill_expirence: jobSkill,
		}

		jobs.totalJobs++
		jobs.ListJobs = append(jobs.ListJobs, job)
	})

	return nil
}

func (jobs *Jobs) GetAllJobs(companies *Companies) error {
	eg := errgroup.Group{}
	if companies.totalCompanies > 0 {
		for i := 0; i < companies.totalCompanies; i++ {
			// https://golang.org/doc/faq#closures_and_goroutines
			url := companies.ListCompanies[i].urlC
			eg.Go(func() error {
				err := jobs.getJobsByUrl(url)
				if err != nil {
					return err
				}
				return nil
			})
		}
		if err := eg.Wait(); err != nil { // Error Group chờ đợi các group goroutines done, nếu có lỗi thì trả về
			return err
		}
	}
	return nil
}

func saveJsonFile(j interface{}, dir string){
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
		} else {
			log.Println(err)
		}
	}

	path := fmt.Sprint(dir, "companies")

	b, err := json.Marshal(j)
	if err != nil { log.Println(err) }

	ioutil.WriteFile(path, b, 0644)
}
