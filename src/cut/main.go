package main

import (
	"fmt"
	"girlImage/src/tool"
	"github.com/disintegration/imaging"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	var path string
	fmt.Printf("请输入图片路径:")
	fmt.Scanln(&path)
	if path == "" {
		path = "/usr/local/www/golang/girlImage/images/meizitu"
	}
	if !tool.Exist(path) {
		fmt.Println("Error:目录不存在")
		tool.Wait()
		return
	}
	defer tool.End()
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && !strings.Contains(path, ".DS_Store") {
			wg.Add(1)
			go cutFile(path)
		}
		return nil
	})
	wg.Wait()
}

func cutFile(oldFile string) {
	fmt.Println(oldFile)
	defer wg.Done()
	newFile := strings.Replace(oldFile, "images", "new-images", 1)
	tool.CheckDir(newFile[0:strings.LastIndex(newFile, "/")])
	src, err := imaging.Open(oldFile)
	if err != nil {
		fmt.Printf("failed to open file: %v, error:%v", oldFile, err)
		return
	}
	src = imaging.Fill(src, src.Bounds().Max.X, src.Bounds().Max.Y-30, imaging.TopLeft, imaging.Lanczos)
	err = imaging.Save(src, newFile)
	if err != nil {
		fmt.Printf("failed to save file: %v, error:%v", oldFile, err)
		return
	}
}
