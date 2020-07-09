package main

import (
	"girlImage/src/tool"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"sync"
)

const URL = "http://jandan.net/ooxx"

var ImgPath = "/images/jiandan"
var Header = map[string]string{
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
}

var dwg sync.WaitGroup
var pwg sync.WaitGroup //分页等待组

func main() {
	dir, _ := os.Getwd()
	ImgPath = dir + ImgPath
	tool.CheckDir(ImgPath)
	pwg.Add(1)
	startDownload(URL)
	pwg.Wait()
	dwg.Wait()
	defer tool.End()
}

func startDownload(url string) {
	resp, err := tool.Get(url, Header)
	defer pwg.Done()
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	nextUrl, exists := doc.Find(".previous-comment-page").Attr("href")
	if exists {
		pwg.Add(1)
		nextUrl = "https:" + nextUrl
		go startDownload(nextUrl)
	}
	var ele = doc.Find(".commentlist").Find("img")
	dwg.Add(ele.Length())
	ele.Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("src")
		go tool.SaveFile("https:"+url, ImgPath, Header, 0, &dwg)
	})
}
