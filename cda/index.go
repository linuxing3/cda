// Command keys is a chromedp example demonstrating how to send key events to
// an element.
package cda

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/linuxing3/cda/config"
)

var htmlContent string


func StartCrawlCda() {

	c := config.InitChromedpOptions(true)
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	defer cancel()

	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//创建一个上下文超时时间为40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 3600*time.Second)
	defer cancel()

	// run task list
	err := chromedp.Run(timeoutCtx, LoginCda(CdaBaseURL), FetchCdaCourseList())
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(htmlContent)
	courses := GetCdaCoursesWithDetails(htmlContent)

	file, _ := json.MarshalIndent(courses, "", " ")
	_ = ioutil.WriteFile("courses.json", file, 0644)

	// err = chromedp.Run(timeoutCtx, 
	// 	loginCda(CdaBaseURL),
	// 	chromedp.Navigate(testLink),
	// 	chromedp.Click(btn, chromedp.NodeVisible),
	// 	chromedp.WaitVisible("video", chromedp.NodeVisible),
	// 	chromedp.Click(confirm, chromedp.NodeVisible),
	// )
	// if err != nil {
	// 	fmt.Println(err)
	// }
}

