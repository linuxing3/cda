package channel

import (
	"fmt"
)

func CrawlWithChannel(url string)  (chan []string, chan string){
	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	data := []string{url}

	// go func() { worklist <- os.Args[1:] }()
	// 初始化协程: 建立工作者列表
	go func() { worklist <- data }()

	// 次协程: 5个爬虫抓取每一个未读过的链接.
	for i := 0; i < 5; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := Crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// 主协程: 将worklist中项目发送
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}

	return worklist, unseenLinks
}

// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func Crawl(url string) []string {
    fmt.Println(url)
    tokens <- struct{}{} // acquire a token
    list:= []string{url}
    <-tokens // release the token
    return list
}