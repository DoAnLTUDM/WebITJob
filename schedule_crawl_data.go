package main

import (
    "github.com/robfig/cron"
    "log"
)

func schedule_crawl_data(){
    c := cron.New()
    c.AddFunc("@daily", func() {
        crawl_data("https://itviec.com/vi")
        crawl_data("https://www.topitworks.com/vi/cong-ty-it-hang-dau-viet-nam")
        log.Println("Daily crawing data")
    })
    c.Start()
}
