package main

import (
	"fmt"
	"girlImage/src/tool"
	"github.com/disintegration/imaging"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var wg sync.WaitGroup

// 保护方式允许一个函数
func ProtectRun(entry func()) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("recover runtime error:", err)
		default: // 非运行时错误
			fmt.Println("recover error:", err)
		}
	}()
	entry()
}

func main() {
	// 故意造成空指针访问错误
	ProtectRun(func() {
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
	})
	tool.Wait()
}

func cutFile(oldFile string) {
	fmt.Println(oldFile)
	defer wg.Done()
	newFile := strings.Replace(oldFile, "images", "new-images", 1)
	char := "/"
	if runtime.GOOS == "windows" {
		char = "\\"
	}
	tool.CheckDir(newFile[0:strings.LastIndex(newFile, char)])
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
