package main

import (
	"fmt"
	"girlImage/src/tool"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const URL = "https://www.mzitu.com/mm/page/"

var ImgPath = "/images/meizitu/"

var Header = map[string]string{
	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36",
	"referer":    URL,
}

var pageIndex = 1
var delay = 500 //每次请求延迟
var wg sync.WaitGroup

func main() {
	dir, _ := os.Getwd()
	ImgPath = dir + ImgPath
	defer tool.End()
	startDownload()
}

func startDownload() {
	currentUrl := URL + strconv.Itoa(pageIndex)
	resp, err := tool.Get(currentUrl, Header)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("#pins").Find("span a").Each(func(i int, s *goquery.Selection) {
		detailUrl, _ := s.Attr("href")
		name := s.Text()
		getDetail(detailUrl, ImgPath+name)
		time.Sleep(time.Millisecond * 100)
	})
	pageIndex++
	startDownload()
}

func getDetail(url string, dirName string) {
	resp, err := tool.Get(url, Header)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	img := doc.Find(".main-image img")
	firstImg, _ := img.Attr("src")
	lastDoc := doc.Find(".pagenavi a")
	lastIndex := lastDoc.Length() - 2
	totalCount := 0
	lastDoc.Each(func(i int, s *goquery.Selection) {
		if i == lastIndex {
			lastUrl, _ := s.Attr("href")
			urlArr := strings.Split(lastUrl, "/")
			totalCount, _ = strconv.Atoi(urlArr[len(urlArr)-1])
		}
	})
	tool.CheckDir(dirName)
	for i := 1; i <= totalCount; i++ {
		imgUrl := strings.Replace(firstImg, "01.jpg", fmt.Sprintf("%02d", i)+".jpg", 1)
		wg.Add(1)
		tool.SaveFile(imgUrl, dirName, Header, delay, &wg)
	}
}
