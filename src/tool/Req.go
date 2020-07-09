package tool

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var tryCount = 0

func Get(url string, header map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	for key, val := range header {
		req.Header.Set(key, val)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		tryCount++
		if tryCount > 3 {
			return nil, err
		}
		fmt.Printf("网络请求失败:%v\r", err)
		time.Sleep(time.Second * 3)
		return Get(url, header)
	}
	if resp.StatusCode != 200 {
		tryCount++
		if tryCount > 3 {
			fmt.Printf("%v ,尝试多次失败: %v\n", url, resp.Status)
			return nil, errors.New("重试多次失败" + resp.Status)
		}
		fmt.Printf("网络请求失败:%s\n3秒后开始重试第%d次...\r", resp.Status, tryCount)
		time.Sleep(time.Second * 3)
		return Get(url, header)
	}
	if tryCount > 0 {
		tryCount = 0
	}
	return resp, nil
}
