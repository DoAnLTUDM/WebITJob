package main

import (
    "log"
    "github.com/PuerkitoBio/goquery"
    "os"
    "strings"
    "regexp"
    "golang.org/x/sync/errgroup"
)


type Job struct {
    Title string
    Salary string
    Address string
    Time_posted string
    Job_reason string
    Job_description string
    Skill_expirence string
    Qualification string
    Company_name string
}

type Jobs struct{
    TotalJobs int
    ListJobs []Job
}

type Company struct{
    Name string
    Address string
    Country string
    Logo string
    UrlC string
}

type Companies struct{
    totalCompanies int
    ListCompanies []Company
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

    jobs := NewJobs()
    err = jobs.GetAllJobs(companies)
    checkError(err)
}

func (companies *Companies) getCompaniesByUrl (urlC string) error{
    doc, err := goquery.NewDocument(urlC)
    if err != nil {
        return err
    }

    doc.Find("#container .top-companies .col-md-4 ").Each(func(i int, selection *goquery.Selection) {
        companyUrl,exist := selection.Find("a").Attr("href")
        if exist{
            if strings.Contains(urlC, "itviec"){
                companyUrl = "https://itviec.com" + companyUrl
            }
        }else {
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
        re := regexp.MustCompile(`\r?\n`)
        companyImg,_ := selection.Find(".headers .logo-container img").Attr("src")
        companyName := re.ReplaceAllString(strings.TrimSpace(selection.Find(
            ".headers .name-and-info .title").Text()), " ")
        var imgName string
        locat := "public/img/company-logo/"
        if strings.Contains(companyImg, "png"){
            tmp := companyName
            imgName = strings.Replace(tmp," ","",-1) + ".png"
        }else{
            tmp := companyName
            imgName = strings.Replace(tmp," ","",-1) + ".jpg"
        }
        download(companyImg, locat, imgName)
        companyAddr := re.ReplaceAllString(strings.TrimSpace(selection.Find(".col-md-3 .map-address").Text())," ")
        companyCountry := selection.Find(".headers .company-info .country span").Text()

        company := Company{
            Name:companyName,
            UrlC:companyUrl,
            Logo:strings.Replace(locat,"public", "", -1)+imgName,
            Address:companyAddr,
            Country:companyCountry,
        }
        company.Create()
        companies.totalCompanies++
        companies.ListCompanies = append(companies.ListCompanies, company)
    })
    return nil
}

func (jobs *Jobs) GetAllJobs(companies *Companies) error {
    eg := errgroup.Group{}
    if companies.totalCompanies > 0 {
        for i := 0; i < companies.totalCompanies; i++ {
            // https://golang.org/doc/faq#closures_and_goroutines
            url := companies.ListCompanies[i].UrlC
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

        job := Job {
            Title:jobTitle,
            Company_name: jobCompany,
            Salary: jobSalary,
            Address: jobAddress,
            Time_posted: jobTimePosted,
            Job_reason: jobReason ,
            Job_description: jobDescription,
            Skill_expirence: jobSkill,
        }
        job.Create()
        jobs.TotalJobs++
        jobs.ListJobs = append(jobs.ListJobs, job)
    })
    return nil
}
