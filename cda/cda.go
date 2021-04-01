package cda

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/linuxing3/cda/config"
	"github.com/spf13/viper"
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
	CdaBaseURL              = "https://e-cda.cn"
	Username                = viper.GetString("config.user.username")
	Password                = viper.GetString("config.user.password")
	CdaLoginFormLoginBtn    = "form > input.login_btn "
	CdaCourseRow            = "#module_0  .hoz_course_row"  // rows
	CdaCourseSelector       = "h2.hoz_course_name  a"       // course name
	CdaCourseProgressBar    = "span.h_pro_percent"          // course progress %
	CdaChooseOpenVideoBtn   = "#module_0 .rt.btn_group > a" // course video link
	CdaChooseConfirmPlayBtn = ".continue > .user_choise"    // confirm playvideo
)

// LoginCda 登录
func LoginCda(host string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(host),
		chromedp.WaitVisible(`.banner_btn`, chromedp.NodeVisible),
		chromedp.Click(".banner_btn", chromedp.NodeVisible),
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.WaitVisible(`#pwd`, chromedp.ByID),
		chromedp.SetValue(`#username`, Username, chromedp.ByID),
		chromedp.SetValue(`#pwd`, Password, chromedp.ByID),
		chromedp.Click("input.login_btn", chromedp.NodeVisible),
	}
}

// FetchCdaCourseList 打开课程列表
func FetchCdaCourseList() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.WaitVisible("h2.join_course_name", chromedp.NodeVisible),
		chromedp.Click("h2.join_course_name > a", chromedp.NodeVisible),
		chromedp.WaitVisible(CdaCourseRow, chromedp.NodeVisible),
		// FIXED 可以抓取某一课程下面全部子课程的整个页面
		chromedp.OuterHTML(config.DocBodySelector, &htmlContent, chromedp.ByJSPath),
	}
}

// PlayCdaVideo 打开课件视频网页, 然后播放视频, 用于模拟用户实际观看
func PlayCdaVideo() chromedp.Tasks {
	return chromedp.Tasks{
		// 3. 打开视频课件网页
		chromedp.Click(CdaChooseOpenVideoBtn, chromedp.NodeVisible),
		// 4. 播放视频
		chromedp.WaitVisible("video", chromedp.NodeVisible),
		chromedp.Click(CdaChooseConfirmPlayBtn, chromedp.NodeVisible),
		// 5. 模拟超时
		// chromedp.WaitVisible("test", chromedp.NodeVisible),
	}
}

// GetCdaCoursesWithDetails 不打开视频网页,直接抓取当前页面的课程列表
func GetCdaCoursesWithDetails(htmlContent string) (courses []Course) {

	courseList, err := config.GetDataList(htmlContent, CdaCourseRow)
	if err != nil {
		log.Fatal(err)
	}

	courseList.Each(func(i int, s *goquery.Selection) {
		// 查找课程具体信息
		item := s.Find(CdaCourseSelector)
		url, _ := item.Attr("href")
		title := item.Text()
		progress := s.Find(CdaCourseProgressBar).Text()
		// 课程编号
		var cid string
		cid, _ = s.Find(".hover_btn").Attr("onclick")
		cid = strings.ReplaceAll(cid, "addUrl(", "")
		cid = strings.ReplaceAll(cid, ")", "")
		// 课程网页
		pageURL := strings.Join([]string{CdaBaseURL, url}, "")
		// 视频网页链接
		vPageLink := strings.Join([]string{
			"https://e-cda.cn/portal/play.do?menu=course&id=", cid},
			"")
		// 视频下载链接
		vDownLink := strings.Join([]string{
			"https://cdn.gwypx.com.cn/course/n",
			cid,
			"/",
			"1.mp4"},
			"")

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
