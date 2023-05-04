package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func main() {
	var pages int
	fmt.Println("请输入要爬取的页数:")
	fmt.Scanf("%d", &pages)

	// 根据页数设置协程数
	ch := make(chan struct{}, pages)
	wg := sync.WaitGroup{}
	for i := 1; i <= pages; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go ImageCrawler(i, ch, &wg)
	}
	wg.Wait()

	fmt.Println("爬取完成。")
}

func ImageCrawler(page int, ch chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		<-ch
		wg.Done()
	}()

	err := DownloadByImageUrls("https://www.bizhizu.cn/shouji/tag-%E7%BE%8E%E5%A5%B3/" + strconv.Itoa(page) + ".html")
	if err != nil {
		fmt.Println(err)
	}
}

func DownloadByImageUrls(url string) error {
	html, err := GetHtml(url)
	if err != nil {
		return err
	}

	// 正则
	re := regexp.MustCompile(`https://uploadfile.bizhizu.cn[^"]+?(\.jpg)`)
	results := re.FindAllStringSubmatch(html, -1)
	// [][]string
	for _, result := range results {
		imageUrl := result[0]
		err = DownloadImage(imageUrl)
		if err != nil {
			return err
		}
	}

	return nil
}

// 根据 url 得到整个 html 信息
func GetHtml(url string) (string, error) {
	var html string
	resp, err := http.Get(url)
	if err != nil {
		return html, errors.New("html http.Get:" + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return html, errors.New("html ioutil.ReadAll:" + err.Error())
	}

	html = string(body)
	return html, nil
}

func DownloadImage(url string) error {
	filename := GetFileName(url)
	err := Download(url, filename)
	if err != nil {
		return err
	}
	return nil
}

func GetFileName(url string) string {
	lastIndex := strings.LastIndex(url, "/")
	fileName := url[lastIndex+1:]
	return fileName
}

func Download(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.New("download http.Get:" + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("download ioutil.ReadAll:" + err.Error())
	}

	path := "./images"
	if _, err = os.Stat(path); err != nil {
		err = os.MkdirAll(path, 0711)
		if err != nil {
			return err
		}
	}

	filename = "./images/" + filename
	err = ioutil.WriteFile(filename, body, 0666)
	if err != nil {
		return err
	}
	return nil
}
