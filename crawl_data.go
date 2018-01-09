package main

import (
    "log"
    "strings"
    "regexp"
    "time"
    "strconv"
    "github.com/tebeka/selenium"
)


type Job struct {
    Id int
    IdComp int
    IdSkills []int
    Skill_name []string
    Title string
    Salary string
    SrcImage string
    Address string
    Time_posted time.Time
    Job_reason []string
    Job_reason_map map[string]interface{}
    Job_description []string
    Job_description_map map[string]interface{}
    Skill_expirence []string
    Skill_expirence_map map[string]interface{}
    ImgComp string
}

type Company struct{
    Id int
    Name string
    Address string
    Country string
    Logo string
    Banner string
    Introduce []string
    Introduce_map map[string]interface{}
    IdSkills []int
    Skill_name []string
    Numjobs int
}

type Skill struct{
    Id int
    Name string
}

func checkError(err error){
    if err != nil{
        log.Println(err)
    }
}

func crawl_data(url string)  {

    // FireFox driver without specific version
    caps := selenium.Capabilities{"browserName": "firefox"}
    wd, _ := selenium.NewRemote(caps, "")
    defer wd.Quit()

    // Get simple playground interface
    wd.Get(url)

    checkError(getDataItViec(wd))
}

func  convertTimePost(day string) time.Time {
    re := regexp.MustCompile("[0-9]+")
    //lấy số trong chuỗi string
    timePosted, _:= strconv.Atoi(re.FindAllString(day, -1)[0])
    return time.Now().AddDate(0, 0, -timePosted)
}

func loginItViec(wd selenium.WebDriver){
    //Tìm và click vào nút đăng nhập
    btn, err := wd.FindElement(selenium.ByXPATH, "(//a[contains(text(),'Đăng Nhập')])[2]")
    if err != nil{
        log.Println(err)
    }
    btn.Click()

    //Nhập vào tên đăng nhập
    email, err := wd.FindElement(selenium.ByXPATH, "(//input[@id='user_email'])[2]")
    if err != nil{
        log.Println(err)
    }
    email.SendKeys("testcast2@gmail.com")

    //Nhập vào mật khẩu
    passwd, err := wd.FindElement(selenium.ByXPATH, "(//input[@id='user_password'])[2]")
    if err != nil{
        log.Println(err)
    }
    passwd.SendKeys("00147Dat")

    //Commit thông tin đăng nhập
    btn, err = wd.FindElement(selenium.ByXPATH, "(//input[@name='commit'])[3]")
    if err != nil{
        log.Println(err)
    }
    btn.Click()
    time.Sleep(time.Second*2)
}

func getDataItViec(wd selenium.WebDriver) error{
    //Mở thẻ "Công ty IT Hàng Đầu"
    if btn, err := wd.FindElement(selenium.ByLinkText, "Công ty IT Hàng Đầu"); err == nil{
        btn.Click()
    }else {
        log.Println(err)
    }

    //Mở tất cả các công ty trong thẻ "Công ty IT Hàng Đầu"
    if btn, err := wd.FindElement(selenium.ByLinkText, "Xem thêm"); err == nil {
        btn.Click()
        time.Sleep(time.Second)
    } else {
        log.Println(err)
    }

    //Đếm số lượng công ty
    count, _ := wd.FindElements(selenium.ByXPATH, "//div[@id='most-popular-companies']/div/div[3]/div/div[@class='col-md-4']")
    for i := 1; i < len(count)+1; i++{
        getTotalInfo(i, wd)
    }
    return nil
}

func getTotalInfo(num int, wd selenium.WebDriver)  {

    companyElement, _ :=	wd.FindElement(selenium.ByXPATH,"//div[@id='most-popular-companies']/div/div[3]/div/div["+strconv.Itoa(num)+"]/a")
    urlCompany, _ := companyElement.GetAttribute("href")

    caps := selenium.Capabilities{"browserName": "firefox"}
    companyWd, _ := selenium.NewRemote(caps, "")
    defer companyWd.Quit()

    companyXPATH := "//div[@class='company-content']/div[@class='company-page']"

    getCompanyInfo(companyWd, urlCompany, companyXPATH)
    //go getJobInfo(countJob, idCompany, jobWd, companyXPATH )
}

func getCompanyInfo(companyWd selenium.WebDriver, urlCompany string, companyXPATH string )  {
    company := Company{}
    // Get simple playground interface
    companyWd.Get("https://itviec.com"+urlCompany)

    getName, _ := companyWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[2]/div[@class='name-and-info']/h1")
    company.Name, _ = getName.Text()
    getBanner, _ := companyWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[1]/img")
    banner,_ := getBanner.GetAttribute("src")
    company.Banner = banner

    getLogo, _ := companyWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[2]/div[@class='logo-container']/div/img")
    logo, _ := getLogo.GetAttribute("src")
    company.Logo = logo

    getCity, _ := companyWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[2]/div[@class='name-and-info']/span")
    company.Address, _ = getCity.Text()

    getCountry, _ := companyWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[2]/div[@class='name-and-info']" +
        "/div[@class='company-info']/div[@class='country']")
    company.Country, _ = getCountry.Text()

    getIntroTitle, _ := companyWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[3]/div/div/div[@class='panel-heading']/h3")
    title, _ := getIntroTitle.Text()
    company.Introduce = append(company.Introduce, title)


    getIntroBody, _ := companyWd.FindElements(selenium.ByXPATH, companyXPATH+"/div[3]/div[1]/div[1]/div[@class='panel-body']" +
    	"/div[1]/p")

    for i:=1; i<len(getIntroBody)+1;i++{
        body, _ := companyWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[3]/div[1]/div[1]/div[@class='panel-body']" +
        	"/div[1]/p["+strconv.Itoa(i)+"]")
        introBody, _ := body.Text()
        company.Introduce = append(company.Introduce, introBody)
    }

    countSkills, _ := companyWd.FindElements(selenium.ByXPATH, companyXPATH+"/div[3]/div/div/" +
        "div[@class='panel-body']/ul[@class='employer-skills']/li")

    for i:=1; i < len(countSkills)+1;i++{
        skill := Skill{}
        getSkill, _ := companyWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[3]/div/div/" +
            "div[@class='panel-body']/ul[@class='employer-skills']/li["+strconv.Itoa(i)+"]")
        skill.Name, _ = getSkill.Text()
        idSkill, _ := skill.Create()
        if idSkill == 0{
            idSkill, _ = getIdSkill(skill.Name)
        }
        company.IdSkills = append(company.IdSkills, idSkill)
    }
    idC, _ := company.Create()
    if idC == 0{
        idC, _ = getIdCompany(company.Name)
    }
    countJobs, _ := companyWd.FindElements(selenium.ByXPATH, companyXPATH+"/div[3]/div[1]/" +
        "div[2]/div[@class='panel-body']/div")

    getJobInfo(len(countJobs), idC, companyWd, companyXPATH)
}

func getJobInfo(countJob int, idCompany int, jobWd selenium.WebDriver, companyXPATH string){
    if(countJob>0){
        loginItViec(jobWd)
        for i:=1; i < countJob+1;i++{
            job := Job{}
            job.IdComp = idCompany
            getTitle, error := jobWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[3]/div[1]/" +
                "div[2]/div[@class='panel-body']/div["+strconv.Itoa(i)+"]/div/div[2]/div/div/h4[@class='title']/a")
            if error == nil{
                job.Title, _ = getTitle.Text()
            }

            getSalary, error:= jobWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[3]/div[1]/" +
                "div[2]/div[@class='panel-body']/div["+strconv.Itoa(i)+"]/div/div[2]/div/div/div[@class='salary']" +
                "/span[@class='salary-text']")
            if error == nil{
                job.Salary, _ = getSalary.Text()
            }

            getTimePost, error := jobWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[3]/div[1]/" +
                "div[2]/div[@class='panel-body']/div["+strconv.Itoa(i)+"]/div/div[2]/div[1]/div[2]/" +
                "div[@class='distance-time-job-posted']/span")
            if error == nil{
                timePost, _ := getTimePost.Text()
                job.Time_posted = convertTimePost(timePost)
            }

            countSkills, error := jobWd.FindElements(selenium.ByXPATH, companyXPATH+"/div[3]/div[1]/" +
                "div[2]/div[@class='panel-body']/div["+strconv.Itoa(i)+"]/div/div[2]/div[2]/div[@class='tag-list']/a")
            if error == nil{
                for j:=1;j<len(countSkills)+1;j++{
                    skill := Skill{}
                    getSkills, _ := jobWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[3]/div[1]/" +
                        "div[2]/div[@class='panel-body']/div["+strconv.Itoa(i)+"]/div/div[2]/div[2]/div[@class='tag-list']" +
                        "/a["+strconv.Itoa(j)+"]/span")
                    skill.Name, _= getSkills.Text()
                    idSkill, _ := skill.Create()
                    if idSkill == 0{
                        idSkill, _ = getIdSkill(skill.Name)
                    }
                    job.IdSkills = append(job.IdSkills, idSkill)
                }
            }
        getUrl, error := jobWd.FindElement(selenium.ByXPATH, companyXPATH+"/div[3]/div[1]/" +
            "div[2]/div[@class='panel-body']/div["+strconv.Itoa(i)+"]/div[@class='job_content']/" +
            "div[@class='job__description']/div[@class='job__body']/div[@class='details']/h4/a")
        if error == nil {
            urlJob, _ := getUrl.GetAttribute("href")
            caps := selenium.Capabilities{"browserName": "firefox"}
            wd, _ := selenium.NewRemote(caps, "")
            defer wd.Quit()
            wd.Get("https://itviec.com"+ urlJob)

            //Get Address
            numAddress,_ := wd.FindElements(selenium.ByXPATH, "//div[@class='job_info']/div[@class='address']" +
                "/div[@class='address__full-address']/span")
            for j:=1;j<len(numAddress)+1;j++{
                getAddress, _ := wd.FindElement(selenium.ByXPATH, "//div[@class='job_info']/div[@class='address']" +
                    "/div[@class='address__full-address']/span["+strconv.Itoa(j)+"]")
                address, _ := getAddress.Text()
                job.Address += address
            }

            //Get Reason
            getTitleReason, _ := wd.FindElement(selenium.ByXPATH, "//div[@class='job-detail']/" +
                "div[@class='job_reason_to_join_us']/h2")
            titleReason, _:= getTitleReason.Text()
            job.Job_reason = append(job.Job_reason, titleReason)
            getReasons, _ := wd.FindElements(selenium.ByXPATH, "//div[@class='job-detail']/" +
                "div[@class='job_reason_to_join_us']/div/ul/li")
            for j:=1;j<len(getReasons)+1;j++{
                getReason, _ := wd.FindElement(selenium.ByXPATH, "//div[@class='job-detail']/" +
                    "div[@class='job_reason_to_join_us']/div/ul/li["+strconv.Itoa(j)+"]")
                reason, _ := getReason.Text()
                job.Job_reason = append(job.Job_reason, reason)
            }

            //Get description
            getTitleDesc, _ := wd.FindElement(selenium.ByXPATH, "//div[@class='job_description']/" +
                "div[@class='title-apply-line']/h2")
            titleDesc, _:= getTitleDesc.Text()
            job.Job_description = append(job.Job_description, titleDesc)
            getDescs, _ := wd.FindElements(selenium.ByXPATH, "//div[@class='job_description']/" +
                "div[@class='description']/ul")
            if len(getDescs) > 0{
                for k:=1;k<len(getDescs)+1;k++{
                    getdescs, _ := wd.FindElements(selenium.ByXPATH, "//div[@class='job_description']/" +
                        "div[@class='description']/ul["+strconv.Itoa(k)+"]"+"/li")
                    for j:=1;j<len(getdescs)+1;j++{
                        getdesc, _ := wd.FindElement(selenium.ByXPATH, "//div[@class='job_description']/" +
                            "div[@class='description']/ul["+strconv.Itoa(k)+"]"+"/li["+strconv.Itoa(j)+"]")
                        desc, _ := getdesc.Text()
                        job.Job_description = append(job.Job_description, desc)
                    }
                }
            }else{
                getDescs, _ := wd.FindElements(selenium.ByXPATH, "//div[@class='job_description']/" +
                    "div[@class='description']/p")
                for j:=1;j<len(getDescs)+1;j++{
                    getDesc, _ := wd.FindElement(selenium.ByXPATH, "//div[@class='job_description']/" +
                        "div[@class='description']/p["+strconv.Itoa(j)+"]")
                    desc, _ := getDesc.Text()
                    job.Job_description = append(job.Job_description, desc)
                }
            }

            //Get skill
            getTitleSkill, _ := wd.FindElement(selenium.ByXPATH, "//div[@class='skills_experience']/h2")
            titleSkill, _:= getTitleSkill.Text()
            job.Skill_expirence = append(job.Job_description, titleSkill)
            getSkills, _ := wd.FindElements(selenium.ByXPATH, "//div[@class='skills_experience']/" +
                "div[@class='experience']/ul")
            if len(getSkills) > 0{
                for k:=1;k<len(getSkills)+1;k++{
                    getskills, _ := wd.FindElements(selenium.ByXPATH, "//div[@class='skills_experience']/" +
                        "div[@class='experience']/ul["+strconv.Itoa(k)+"]"+"/li")
                    for j:=1;j<len(getskills)+1;j++{
                        getskill, _ := wd.FindElement(selenium.ByXPATH, "//div[@class='skills_experience']/" +
                            "div[@class='experience']/ul["+strconv.Itoa(k)+"]"+"/li["+strconv.Itoa(j)+"]")
                        skill, _ := getskill.Text()
                        job.Job_description = append(job.Job_description, skill)
                    }
                }
            }else{
                getSkills, _ := wd.FindElements(selenium.ByXPATH, "//div[@class='skills_experience']/" +
                    "div[@class='experience']/p")
                for j:=1;j<len(getSkills)+1;j++{
                    getSkill, _ := wd.FindElement(selenium.ByXPATH, "//div[@class='skills_experience']/" +
                        "div[@class='experience']/p["+strconv.Itoa(j)+"]")
                    skill, _ := getSkill.Text()
                    job.Job_description = append(job.Job_description, skill)
                }
            }
        }else {
            log.Println(error)
        }
        job.Create()
        }
    }
}

func checkFormatImg(src string, name string) string{
    if strings.Contains(src, "png"){
        return strings.Replace(name," ","",-1) + ".png"
    }else{
        return strings.Replace(name," ","",-1) + ".jpg"
    }
}