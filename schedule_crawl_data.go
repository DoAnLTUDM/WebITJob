package main

import (
	"github.com/robfig/cron"
	"log"
)

func schedule_crawl_data(){
	c := cron.New()
	c.AddFunc("@daily", func() {
		crawl_data("https://itviec.com/vi")
		crawl_data("https://www.vietnamworks.com/viec-lam-it-phan-cung-mang-i55-vn")
		log.Println("Daily crawing data")
	})
	c.Start()
}
