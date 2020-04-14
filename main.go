package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const URL = "http://jandan.net/ooxx"
const IMG_PATH = "images"

var success_count, exists_count, error_count = 0, 0, 0

var dwg sync.WaitGroup //下载等待组
var pwg sync.WaitGroup //分页等待组

var rwMutex *sync.RWMutex

func main() {
	t := time.Now()
	rwMutex = new(sync.RWMutex)
	pwg.Add(1)
	startDownload(URL)
	pwg.Wait()
	dwg.Wait()
	elapsed := time.Since(t)
	fmt.Println("Done\nElapsed time:", elapsed)
	fmt.Print("Enter any key to exit the program:")
	var enter string
	fmt.Scanln(&enter)
}

func startDownload(url string) {
	defer pwg.Done()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.92 Safari/537.36")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal("网络请求失败:", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("网络请求失败: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	nextUrl, exists := doc.Find(".previous-comment-page").Attr("href")
	if (exists) {
		nextUrl = "https:" + nextUrl
		pwg.Add(1)
		go startDownload(nextUrl)
	}

	doc.Find(".commentlist").Find("img").Each(func(i int, s *goquery.Selection) {
		dwg.Add(1)
		url, _ := s.Attr("src")
		go saveFile("https:" + url)
	})

}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func getFileName(url string) string {
	url_list := strings.Split(url, "/")
	return IMG_PATH + "/" + url_list[len(url_list)-1]
}

func saveFile(url string) string {
	defer dwg.Done()
	defer printLog()

	if !exist(IMG_PATH) {
		_ = os.Mkdir(IMG_PATH, 0777)
	}
	filename := getFileName(url)
	if exist(filename) {
		calcCount(&exists_count)
		return ""
	}
	resp, err := http.Get(url)
	if err != nil {
		calcCount(&error_count)
		fmt.Println("失败:"+url, err)
		return ""
	}
	defer resp.Body.Close()
	pix, _ := ioutil.ReadAll(resp.Body)

	if err := ioutil.WriteFile(filename, pix, 0777); err != nil {
		calcCount(&error_count)
		fmt.Println("失败:"+url, err)
		return ""
	}
	calcCount(&success_count)
	return filename
}

func calcCount(t *int) {
	rwMutex.Lock()
	defer rwMutex.Unlock()
	*t = *t + 1
}

func printLog() {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	fmt.Printf("Successed:%v, Existed:%v, Failure:%v\n", success_count, exists_count, error_count)
}
