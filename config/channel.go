package channel

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

// Crawl 抓取一个字符串, 改成一个数组
func Crawl(url string) []string {
	fmt.Println(url)
	// tokens 作为计数信号量, 控制2请求的限制
	var tokens = make(chan struct{}, 20)
	tokens <- struct{}{} // 抓取

	// 将url放进数组里传回
	list := []string{url}

	<-tokens // 释放

	return list
}

// TransformStringToStringArray 抓取unseenLink, 从字符串改成一个数组, 写入worklist通道
func TransformStringToStringArray(worklist chan []string, unseenLinks chan string) {
	for link := range unseenLinks {
		foundLinks := []string{ link }
		go func() { worklist <- foundLinks }()
	}
}

// ClassifyItems 将worklist中项目取出, 判断是否已读, 未读链接发送给未读链接
func ClassifyItems(worklist chan []string, unseenLinks chan string) {
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

func CrawlWithChannel(url string) {
	worklist := make(chan []string)  // URLs,可能有重复
	unseenLinks := make(chan string) // 去重 URLs

	// go func() { worklist <- os.Args[1:] }()
	// 协程1: worklist中得到1个元素
	go func() { worklist <- []string{url} }()

	// 协程2: 5个爬虫
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go TransformStringToStringArray(worklist, unseenLinks)
	}

	// 协程3: 完成去重
	ClassifyItems(worklist, unseenLinks)

	wg.Wait()
}

func SendInt(c chan int, v int) {
	defer wg.Done()
	c <- v
}

func SendString(c chan string, s string) {
	defer wg.Done()
	c <- s
}
