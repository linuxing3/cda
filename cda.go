package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Course struct {
	Title     string
	Url       string
	Progress  string
	Cid       string
	VPageLink string
	VDownLink string
}

var (
	CdaBaseURL            = "https://e-cda.cn"
	CdaLoginFormLoginBtn  = "form > input.login_btn "
	CdaCourseRow          = "#module_0  .hoz_course_row" // rows
	CdaCourseSelector     = "h2.hoz_course_name  a"      // course name
	CdaCourseProgressBar  = "span.h_pro_percent"         // course progress %
	CdaChooseCourseBtn    = ".rt.btn_group > a"          // course video link
	CdaChooseVideoConfirm = ".continue > .user_choise"   // confirm playvideo
)

// sendkeys sends keys to the server and extracts 4 values from the html page.
func LoginCda(host string) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(host),
		chromedp.WaitVisible(`.banner_btn`, chromedp.NodeVisible),
		chromedp.Click(".banner_btn", chromedp.NodeVisible),
		// 1. 登录
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.WaitVisible(`#pwd`, chromedp.ByID),
		chromedp.SetValue(`#username`, "xing_wenju@mfa.gov.cn", chromedp.ByID),
		chromedp.SetValue(`#pwd`, "Abcd1234", chromedp.ByID),
		chromedp.Click("input.login_btn", chromedp.NodeVisible),
	}
}

func FetchCdaCourseList() chromedp.Tasks {

	// btn := "#module_0 .rt.btn_group > a"
	// confirm := ".continue > .user_choise"

	return chromedp.Tasks{
		// 2. 打开课程列表
		chromedp.WaitVisible("h2.join_course_name", chromedp.NodeVisible),
		// div.join_special_main > div > div:nth-child(1) > div.clearfix > ul > li:nth-child(1) > h2
		chromedp.Click("h2.join_course_name > a", chromedp.NodeVisible),

		chromedp.WaitVisible(CdaCourseRow, chromedp.NodeVisible),

		// FIXED 可以抓取某一课程下面全部子课程的整个页面
		chromedp.OuterHTML(DocBodySelector, &htmlContent, chromedp.ByJSPath),
	}
}

func PlayCdaVideo() chromedp.Tasks {

	return chromedp.Tasks{

		// 3. 打开视频列表
		chromedp.Click(btn, chromedp.NodeVisible),

		// // 4. 视频播放
		chromedp.WaitVisible("video", chromedp.NodeVisible),
		chromedp.Click(confirm, chromedp.NodeVisible),

		// // 模拟超时
		// chromedp.WaitVisible("test", chromedp.NodeVisible),
	}
}

// GetCdaCoursesWithDetails 不打开视频网页,而是抓取当前页面的课程列表
func GetCdaCoursesWithDetails(htmlContent string) (courses []Course) {

	courseList, err := GetDataList(htmlContent, CdaCourseRow)
	if err != nil {
		log.Fatal(err)
	}

	courseList.Each(func(i int, s *goquery.Selection) {
		// 查找课程具体信息
		item := s.Find("ul h2.hoz_course_name a")
		url, _ := item.Attr("href")
		title := item.Text()
		progress := s.Find(".h_pro_percent").Text()

		var cid string

		cid, _ = s.Find(".hover_btn").Attr("onclick")
		cid = strings.ReplaceAll(cid, "addUrl(", "")
		cid = strings.ReplaceAll(cid, ")", "")

		pageURL := strings.Join([]string{CdaBaseURL, url}, "")

		vPageLink := strings.Join([]string{"https://e-cda.cn/portal/play.do?menu=course&id=", cid}, "")
		vDownLink := strings.Join([]string{"https://cdn.gwypx.com.cn/course/n", cid, "/", "1.mp4"}, "")

		// videoURL := strings.Join([]string{CdaBaseURL, "/student/class_detail_course.do?", fmt.Sprintf("cid=%s", cid), "&elective_type=1&menu=myclass&tab_index=0"}, "")

		courses = append(courses, Course{
			Title:     title,
			Url:       pageURL,
			Progress:  progress,
			Cid:       cid,
			VPageLink: vPageLink,
			VDownLink: vDownLink,
		})
		// fmt.Println(title)
	})
	return courses
}
