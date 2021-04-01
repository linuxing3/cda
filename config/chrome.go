package config

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

var (
	DocBodySelector = "document.querySelector('body')"
	webSite         = ""
	keyWord         = ""
	ChromeAgentConfig = `Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`
)

// GetHTTPHtmlContent 获取网站上爬取的数据
// url [string] 网址
// selector [string] 必须显示的元素
// selected [interface] 要抓取的元素
func GetHTTPHtmlContent(url string, selector string, selected interface{}) (htmlContent string, err error) {
	c := InitChromedpOptions(true)
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	defer cancel()

	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//创建一个上下文, 超时时间为600s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 600*time.Second)
	defer cancel()

	if err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(selected, &htmlContent, chromedp.ByJSPath),
	); err != nil {
		return "", err
	}
	return htmlContent, nil
}


// GetDataList 得到数据列表
func GetDataList(htmlContent string, selector string) (*goquery.Selection, error) {
	if dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent)); err != nil {
		return nil, err
	} else {
		list := dom.Find(selector)
		return list, nil
	}
}

// InitChromedpOptions 控制是否以无头方式启用Chrome, 初始化参数先传一个空的数据 
func InitChromedpOptions(headless bool) context.Context {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", headless), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(ChromeAgentConfig),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), options...)
	return ctx
}
