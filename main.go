package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	defer pw.Stop()
	pwConfig := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	}
	browser, err := pw.Chromium.Launch(pwConfig)
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	defer browser.Close()
	// Create a new browser context
	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("could not create context: %v", err)
	}

	// FacebookPOC(context)
	// TwitterPOC(context)
	// InstagramPOC(context)
	// YoutubePOC(context)
	TiktokPOC(context)
}

func TwitterPOC(context playwright.BrowserContext) {
	cookies := []playwright.OptionalCookie{
		{
			Name:     "guest_id",
			Value:    "v1%3A170123350582410572",
			Domain:   playwright.String(".x.com"),
			Path:     playwright.String("/"),
			Secure:   playwright.Bool(true),
			HttpOnly: playwright.Bool(true),
			Expires:  playwright.Float(1753846024),
		},
		{
			Name:     "twid",
			Value:    "u%3D1578387547251650560",
			Domain:   playwright.String(".x.com"),
			Path:     playwright.String("/"),
			Secure:   playwright.Bool(true),
			HttpOnly: playwright.Bool(true),
			Expires:  playwright.Float(1753926911),
		},
		{
			Name:     "auth_token",
			Value:    "917218e5d484930bb3d65fe1a268c096f2dd8b06",
			Domain:   playwright.String(".x.com"),
			Path:     playwright.String("/"),
			Secure:   playwright.Bool(true),
			HttpOnly: playwright.Bool(true),
			Expires:  playwright.Float(1753846024),
		},
		{
			Name:     "guest_id_ads",
			Value:    "v1%3A170123350582410572",
			Domain:   playwright.String(".x.com"),
			Path:     playwright.String("/"),
			Secure:   playwright.Bool(true),
			HttpOnly: playwright.Bool(true),
			Expires:  playwright.Float(1756950911),
		},
		{
			Name:     "guest_id_marketing",
			Value:    "v1%3A170123350582410572",
			Domain:   playwright.String(".x.com"),
			Path:     playwright.String("/"),
			Secure:   playwright.Bool(true),
			HttpOnly: playwright.Bool(true),
			Expires:  playwright.Float(1756950911),
		},
		{
			Name:     "ct0",
			Value:    "91f9bf086bd6e729ed2781b772010aa0e6ca4f18d63bc577d1a534fd3713c652f926be20fc5f882e7cc36b01c2b55d9df2863851945142619bf8514265dbf56b0bcc8d9624a85b0dcf8982543fa43039",
			Domain:   playwright.String(".x.com"),
			Path:     playwright.String("/"),
			Secure:   playwright.Bool(true),
			HttpOnly: playwright.Bool(true),
			Expires:  playwright.Float(1756870025),
		},
		{
			Name:     "lang",
			Value:    "en",
			Domain:   playwright.String(".x.com"),
			Path:     playwright.String("/"),
			Secure:   playwright.Bool(true),
			HttpOnly: playwright.Bool(true),
		},
		{
			Name:     "external_referer",
			Value:    "padhuUp37zjgzgv1mFWxJ12Ozwit7owX|0|8e8t2xd8A2w%3D",
			Domain:   playwright.String(".x.com"),
			Path:     playwright.String("/"),
			Secure:   playwright.Bool(true),
			HttpOnly: playwright.Bool(true),
			Expires:  playwright.Float(1722995708),
		},
		{
			Name:     "personalization_id",
			Value:    "\"v1_dgj0eitx/8R49kNkyVpmXQ==\"",
			Domain:   playwright.String(".x.com"),
			Path:     playwright.String("/"),
			Secure:   playwright.Bool(true),
			HttpOnly: playwright.Bool(true),
			Expires:  playwright.Float(1756950911),
		},
	}
	err := context.AddCookies(cookies)
	if err != nil {
		log.Fatalf("could not set cookies: %v", err)
	}
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	if _, err := page.Goto("https://x.com/kem_atrbpn"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	time.Sleep(time.Second * 10)
	articlesLocator := page.Locator(`article[data-testid="tweet"]`)
	postCount, _ := articlesLocator.Count()
	fmt.Println("Total post: ", postCount)

	// Write CSV data to a file
	file, err := os.OpenFile("twitter.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer file.Close()

	csvData := "username,post,date,url,retweet,is_repost\n"
	_, err = file.WriteString(csvData)
	if err != nil {
		log.Fatalf("could not write to file: %v", err)
	}

	for i := 0; i < postCount; i++ {
		articleLocator := articlesLocator.Nth(i)
		content, _ := articleLocator.Locator(`div[data-testid="tweetText"]`).InnerText()
		fmt.Println("POST: ", content[:10])

		headLocator := articleLocator.Locator(`div.css-175oi2r.r-18u37iz.r-1q142lx`)
		url, _ := headLocator.Locator("a").GetAttribute("href")
		fmt.Println("url: ", url)
		timestamp, _ := headLocator.Locator("time").GetAttribute("datetime")
		fmt.Println("timestamp: ", timestamp)
		retweetLocator, _ := articleLocator.Locator(`button[data-testid="retweet"]`).InnerText()
		fmt.Println("retweet: ", retweetLocator)
		repost := "false"
		if strings.Contains(content, "Repost") {
			repost = "true"
		}
		csvData = fmt.Sprintf("%s, %s,%s,%s,%s,%s\n", "kem_atrbpn", content[:10], timestamp, url, retweetLocator, repost)
		_, err = file.WriteString(csvData)
		if err != nil {
			log.Fatalf("could not write to file: %v", err)
		}
	}
}

func FacebookPOC(context playwright.BrowserContext) {
	// Set cookies
	cookies := []playwright.OptionalCookie{
		{
			Name:     "sb",
			Value:    "w8wKZTot5Z9xPpya4x_fwFMN",
			Domain:   playwright.String(".facebook.com"),
			Path:     playwright.String("/"),
			Expires:  playwright.Float(1735792837),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
		},
		{
			Name:     "datr",
			Value:    "w8wKZfi5sv3aYuDFioLGSBtj",
			Domain:   playwright.String(".facebook.com"),
			Path:     playwright.String("/"),
			Expires:  playwright.Float(1735792829),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
		},
		{
			Name:     "c_user",
			Value:    "100002987529086",
			Domain:   playwright.String(".facebook.com"),
			Path:     playwright.String("/"),
			Expires:  playwright.Float(1753840787),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
		},
		{
			Name:     "xs",
			Value:    "27%3AEKfvknscp5mEiw%3A2%3A1701232833%3A-1%3A11172%3A%3AAcVTsgblXiiTcB4OXC7p51Y4ySHq-EZVD5Pvr9P5ytU",
			Domain:   playwright.String(".facebook.com"),
			Path:     playwright.String("/"),
			Expires:  playwright.Float(1753840787),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
		},
		{
			Name:     "wd",
			Value:    "1728x993",
			Domain:   playwright.String(".facebook.com"),
			Path:     playwright.String("/"),
			Expires:  playwright.Float(1722909660),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
		},
		{
			Name:     "fr",
			Value:    "15JXJQgwYjwrv5Zbd.AWWOFWIze3RuBrxGjTGLPEc1ffQ.BmqEkR..AAA.0.0.BmqElZ.AWVqrTFsvnM",
			Domain:   playwright.String(".facebook.com"),
			Path:     playwright.String("/"),
			Expires:  playwright.Float(1730080858),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
		},
		{
			Name:     "presence",
			Value:    "C%7B%22t3%22%3A%5B%5D%2C%22utc3%22%3A1722304860319%2C%22v%22%3A1%7D",
			Domain:   playwright.String(".facebook.com"),
			Path:     playwright.String("/"),
			Expires:  playwright.Float(0),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
		},
	}

	err := context.AddCookies(cookies)
	if err != nil {
		log.Fatalf("could not set cookies: %v", err)
	}
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	if _, err := page.Goto("https://mbasic.facebook.com/KementerianATRBPN"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	page.GetByRole("link", playwright.PageGetByRoleOptions{Name: "Timeline"}).Click()
	targetPage := 5
	// Convert results to CSV

	// Write CSV data to a file
	file, err := os.OpenFile("facebook.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer file.Close()

	csvData := "username,post,date,url,like,comment,is_repost\n"
	_, err = file.WriteString(csvData)
	if err != nil {
		log.Fatalf("could not write to file: %v", err)
	}

	results := []map[string]interface{}{}

	for i := 0; i < targetPage; i++ {
		page.WaitForSelector("div .feed", playwright.PageWaitForSelectorOptions{Timeout: playwright.Float(60000)})
		articlesLocator := page.Locator(`article[data-ft='{"tn":"-R"}']`)
		postCount, _ := articlesLocator.Count()
		log.Printf("Number of posts found: %d", postCount)

		// // Extract post details
		for i := 0; i < postCount; i++ {
			articleLocator := articlesLocator.Nth(i)
			content, _ := articleLocator.Locator(`div[data-ft='{"tn":"*s"}']`).InnerText()
			fmt.Println("POST: ", content)
			date, _ := articleLocator.Locator("abbr").InnerText()
			fmt.Println("DATE: ", date)
			link := articleLocator.Locator(`a:text("Full Story")`)
			href, err := link.GetAttribute("href")
			if err != nil {
				log.Fatalf("could not get href attribute: %v", err)
			}
			fmt.Println("URL: https://mbasic.facebook.com", href)
			likeSpan, err := articleLocator.Locator(`span[id^="like_"]`).InnerText()
			if err != nil {
				log.Fatalf("could not find the span with id starting with 'like_': %v", err)
			}
			fmt.Println("like: ", likeSpan)
			commentsLocator := articleLocator.Locator(`a:text-matches("Comments$")`)
			commentsText, _ := commentsLocator.InnerText()
			fmt.Println("Comments: ", commentsText)
			result := map[string]interface{}{
				"post":    content[0:10],
				"date":    date,
				"url":     "https://mbasic.facebook.com" + href,
				"like":    likeSpan,
				"comment": commentsText,
			}
			results = append(results, result)

			repost := "false"
			if strings.Contains(content, "Repost") {
				repost = "true"
			}
			csvData = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n", "KementerianATRBPN", content[:10], date, "https://mbasic.facebook.com"+href, likeSpan, commentsText, repost)
			_, err = file.WriteString(csvData)
			if err != nil {
				log.Fatalf("could not write to file: %v", err)
			}
		}
		page.GetByRole("link", playwright.PageGetByRoleOptions{Name: "See more stories"}).Click()
	}
}

type Result struct {
	Link         string
	TotalComment int
	TotalLike    int
	Username     string
	Date         time.Time
	Post         string
}

func (res Result) ToCSV() string {
	return fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%s\n",
		url.QueryEscape(res.Username),
		url.QueryEscape(res.Post),
		url.QueryEscape(res.Date.String()),
		url.QueryEscape(res.Link),
		strconv.Itoa(res.TotalLike),
		strconv.Itoa(res.TotalComment),
		"false",
	)
}

func InstagramPOC(ctx playwright.BrowserContext) {
	dat, err := os.ReadFile("./cookies-ig-icjonoss.txt")
	if err != nil {
		log.Fatalf("failed to read cookies file: %v", err)
		return
	}

	cookiesRaw := string(dat)
	re := regexp.MustCompile("\t|\n")
	cookieSplitRaw := re.Split(cookiesRaw, -1)
	// parsing empty
	cookieSplit := []string{}
	for _, v := range cookieSplitRaw {
		if len(v) > 0 {
			cookieSplit = append(cookieSplit, v)
		}
	}

	cookies := []playwright.OptionalCookie{}
	for i := 0; i < len(cookieSplit)/7; i++ {
		idx := i * 7
		expire, _ := strconv.ParseFloat(cookieSplit[idx+4], 64)
		nc := playwright.OptionalCookie{
			Name:     cookieSplit[idx+5],
			Value:    cookieSplit[idx+6],
			Domain:   playwright.String(cookieSplit[idx+0]),
			Path:     playwright.String(cookieSplit[idx+2]),
			Expires:  playwright.Float(expire),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
		}
		cookies = append(cookies, nc)
	}

	if err := ctx.AddCookies(cookies); err != nil {
		log.Fatalf("failed to set cookie: %v", err)
	}

	page, err := ctx.NewPage()
	targetURI := "https://www.instagram.com/kementerian.atrbpn/"
	if _, err := page.Goto(targetURI); err != nil {
		log.Fatalf("failed to open page %v : %v", targetURI, err)
	}

	rows, err := page.Locator("main > div > div:nth-child(3) > div > div").All()
	if err != nil {
		log.Fatalf("failed to get content rows")
	}

	// helper function
	parseTotal := func(v string) int {
		regex := regexp.MustCompile("\\D")
		str := regex.ReplaceAllString(v, "")
		fmt.Println(str)
		res, _ := strconv.Atoi(str)
		return res
	}

	// Write CSV data to a file
	file, err := os.OpenFile("instagram.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer file.Close()

	csvData := "sep=,\nusername,post,date,url,like,comment,is_repost\n"
	_, err = file.WriteString(csvData)
	if err != nil {
		log.Fatalf("could not write to file: %v", err)
	}

	maxFetchedData := 10
	fetchedData := 0
	results := []Result{}
	// TODO handling infinite scrolling
	// https://stackoverflow.com/questions/69183922/playwright-auto-scroll-to-bottom-of-infinite-scroll-pag
	for _, row := range rows {
		if fetchedData > maxFetchedData {
			break
		}
		columns, err := row.Locator("> div").All()
		if err != nil {
			log.Fatalf("failed to get columns row")
			return
		}

		for _, column := range columns {
			fetchedData++
			if fetchedData > maxFetchedData {
				break
			}
			column.Hover()
			newEntry := Result{
				Link:         "",
				TotalComment: 0,
				TotalLike:    0,
				Username:     "",
				Date:         time.Time{},
				Post:         "",
			}

			onHoverContent, err := column.Locator("a").AllInnerTexts()
			if err != nil {
				log.Fatalf("failed to get like & comment: %v", err)
			}
			contentSplit := strings.Split(onHoverContent[0], "\n")
			commentTotal := parseTotal(contentSplit[0])
			likeTotal := parseTotal(contentSplit[1])
			fmt.Printf("comment total: %v, like total: %v \n", commentTotal, likeTotal)
			newEntry.TotalComment = commentTotal
			newEntry.TotalLike = likeTotal

			link, err := column.Locator("a").GetAttribute("href")
			if err != nil {
				log.Fatalf("failed to get link: %v", err)
			}
			fmt.Printf("post link: %v \n", link)
			newEntry.Link = fmt.Sprintf("https://instagram.com/%v", link)

			column.Click()
			postDescLocator := "article[role=\"presentation\"] > div > div:nth-child(2) > div > div > div:nth-child(2) > div > ul > div:first-child"
			descLocator := page.Locator(postDescLocator).First()

			selectorUsername := "> li > div > div > div:nth-child(2) > h2 a"
			if err != nil {
				log.Fatalf("failed to get username: %v", err)
			}
			username, err := descLocator.Locator(selectorUsername).InnerText()
			fmt.Printf("username: %v \n", username)
			newEntry.Username = username

			descTextLocator := "> li > div > div > div:nth-child(2) > div"
			if err != nil {
				log.Fatalf("failed to get description text: %v", err)
			}
			descText, err := descLocator.Locator(descTextLocator).Nth(0).InnerText()
			fmt.Printf("desc: %v \n", descText)
			newEntry.Post = descText

			rawDateTime, err := descLocator.Locator(descTextLocator).Nth(1).Locator("> span time[datetime]").GetAttribute("datetime")
			if err != nil {
				log.Fatalf("failed to get datetime: %v", err)
			}
			dateTime, err := time.Parse("2006-01-02T15:04:05.000Z", rawDateTime)
			if err != nil {
				log.Fatalf("failed to parse time: %v", err)
			}
			fmt.Printf("desc: %v \n", dateTime)
			newEntry.Date = dateTime

			if err := page.Locator("svg[aria-label=\"Close\"]").First().Click(); err != nil {
				log.Fatalf("failed to close post modal: %v", err)
			}
			if err := page.WaitForURL(targetURI); err != nil {
				log.Fatalf("browser is not redirected back to profile page: %v", err)
			}
			results = append(results, newEntry)

			csvData = newEntry.ToCSV()
			_, err = file.WriteString(csvData)
			if err != nil {
				log.Fatalf("could not write to file: %v", err)
			} else {
				fmt.Printf("post dumped to Csv: %v", fetchedData)
			}
		}
	}
}

type ResultYoutube struct {
	Link         string
	TotalComment int
	TotalLike    int
	Username     string
	Date         time.Time
	Desc         string
	Views        int
	Title        string
}

func (res ResultYoutube) ToCSV() string {
	return fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%s,%s\n",
		url.QueryEscape(res.Username),
		url.QueryEscape(res.Title),
		url.QueryEscape(res.Desc),
		url.QueryEscape(res.Date.String()),
		url.QueryEscape(res.Link),
		strconv.Itoa(res.TotalLike),
		strconv.Itoa(res.TotalComment),
		strconv.Itoa(res.Views),
	)
}

func YoutubePOC(ctx playwright.BrowserContext) {
	page, err := ctx.NewPage()
	if err != nil {
		log.Fatalf("failed to open new page: %v", err)
	}
	targetURI := "https://www.youtube.com/@KementerianATRBPN/videos"
	if _, err := page.Goto(targetURI); err != nil {
		log.Fatalf("failed to go to %v: %v", targetURI, err)
	}

	contents, err := page.Locator("ytd-rich-item-renderer").All()
	if err != nil {
		log.Fatalf("failed to get contents wrapper: %v", err)
	}
	maxFetchedData := 10
	fetchedData := 0

	results := []ResultYoutube{}
	// Write CSV data to a file
	file, err := os.OpenFile("youtube.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer file.Close()

	csvData := fmt.Sprintf("sep=,\n%v\n", "username,title,desc,date,link,like,comment,views")
	_, err = file.WriteString(csvData)
	if err != nil {
		log.Fatalf("could not write to file: %v", err)
	}

	for _, content := range contents {
		if fetchedData >= maxFetchedData {
			break
		}
		newEntry := ResultYoutube{}
		fetchedData++
		vidURL, _ := content.Locator("a").First().GetAttribute("href")
		newEntry.Link = fmt.Sprintf("https://youtube.com/%v", vidURL)
		newEntry.Username = "@KementerianATRBPN"
		results = append(results, newEntry)
	}

	for i, res := range results {
		if _, err := page.Goto(res.Link); err != nil {
			log.Fatalf("failed to open vid link %v: %v", res.Link, err)
		}

		if err := page.Locator("div#title h1").First().WaitFor(); err != nil {
			log.Fatalf("failed to load async video data: %v", err)
		}

		results[i].Title, err = page.Locator("div#title h1").First().InnerText()
		if err != nil {
			log.Fatalf("failed to get video title: %v", err)
		}

		tooltipValue, err := page.Locator("#description-inner #tooltip").First().InnerText()
		tooltipValue = strings.TrimSpace(tooltipValue)
		if err != nil {
			log.Fatalf("failed to get tooltip value: %v", err)
		}
		regexDate := regexp.MustCompile("\\w\\w\\w \\d?\\d, \\d\\d\\d\\d")
		match := regexDate.FindAllString(tooltipValue, 1)
		if len(match) < 1 {
			log.Fatalf("failed to get upload date: tooltipvalue = %v", tooltipValue)
		}
		if results[i].Date, err = time.Parse("Jan _2, 2006", match[0]); err != nil {
			log.Fatalf("failed to parse date %v: %v", match[0], err)
		}

		regexView := regexp.MustCompile("^(.+) view")
		match = regexView.FindStringSubmatch(tooltipValue)
		if len(match) < 2 {
			log.Fatalf("invalid regex results for total views %v", tooltipValue)
		}
		results[i].Views, err = strconv.Atoi(strings.ReplaceAll(match[1], ",", ""))
		if err != nil {
			log.Fatalf("failed to parse total views %v: %v", match[1], err)
		}

		// just in case
		page.Locator("#description-inline-expander #expand").First().Click()
		page.Locator("#description-inner #description-inline-expander yt-attributed-string").First().WaitFor()
		results[i].Desc, _ = page.Locator("#description-inner #description-inline-expander yt-attributed-string").First().InnerText()
		// skip error check, since maybe the video didn't have description

		btnLikeLabel, err := page.Locator("like-button-view-model button").First().GetAttribute("aria-label")
		if err != nil {
			log.Fatalf("failed to get like button: %v", btnLikeLabel)
		}
		regexLikeCount := regexp.MustCompile("(\\d|,)+")
		match = regexLikeCount.FindStringSubmatch(btnLikeLabel)
		if len(match) < 1 {
			log.Fatalf("invalid regex results for total likes %v", btnLikeLabel)
		}
		results[i].TotalLike, err = strconv.Atoi(strings.ReplaceAll(match[0], ",", ""))
		if err != nil {
			log.Fatalf("failed to parse total likes %v: %v", match[0], err)
		}
		page.Mouse().Wheel(0, 1200)
		page.Locator("ytd-comments-header-renderer #count span").First().WaitFor()
		totalComment, err := page.Locator("ytd-comments-header-renderer #count span").First().InnerText()
		if totalComment != "" {
			results[i].TotalComment, err = strconv.Atoi(strings.ReplaceAll(totalComment, ",", ""))
			if err != nil {
				log.Fatalf("failed to parse total comments %v: %v", totalComment, err)

			}
		}

		j, _ := json.MarshalIndent(results[i], "", "  ")
		fmt.Println(string(j))

		_, err = file.WriteString(results[i].ToCSV())
		if err != nil {
			log.Fatalf("could not write to file: %v", err)
		}
	}
}

func TiktokPOC(ctx playwright.BrowserContext) {
	dat, err := os.ReadFile("./cookies-tiktok-doidoidoi06.txt")
	if err != nil {
		log.Fatalf("failed to read cookies file: %v", err)
		return
	}

	cookiesRaw := string(dat)
	re := regexp.MustCompile("\t|\n")
	cookieSplitRaw := re.Split(cookiesRaw, -1)
	// parsing empty
	cookieSplit := []string{}
	for _, v := range cookieSplitRaw {
		if len(v) > 0 {
			cookieSplit = append(cookieSplit, v)
		}
	}

	cookies := []playwright.OptionalCookie{}
	for i := 0; i < len(cookieSplit)/7; i++ {
		idx := i * 7
		expire, _ := strconv.ParseFloat(cookieSplit[idx+4], 64)
		nc := playwright.OptionalCookie{
			Name:     cookieSplit[idx+5],
			Value:    cookieSplit[idx+6],
			Domain:   playwright.String(cookieSplit[idx+0]),
			Path:     playwright.String(cookieSplit[idx+2]),
			Expires:  playwright.Float(expire),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
		}
		cookies = append(cookies, nc)
	}

	if err := ctx.AddCookies(cookies); err != nil {
		log.Fatalf("failed to set cookie: %v", err)
	}

	page, err := ctx.NewPage()
	targetURI := "https://www.tiktok.com"
	if _, err := page.Goto(targetURI); err != nil {
		log.Fatalf("failed to open page %v : %v", targetURI, err)
	}

	targetAccount := "@kementerian.atrbpn"
	page.Locator("form[data-e2e='search-box'] input").First().Fill(targetAccount)

	time.Sleep(10 * time.Minute)
}
