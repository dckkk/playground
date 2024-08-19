package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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
	tiktokPOC(context)
}

func tiktokPOC(ctx playwright.BrowserContext) {
	dat, err := os.ReadFile("./www.tiktok.com_cookies.txt")
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
	targetURI := "https://www.tiktok.com/@kanwilbpnkaltim"
	if _, err := page.Goto(targetURI); err != nil {
		log.Fatalf("failed to open page %v : %v", targetURI, err)
	}

	time.Sleep(10 * time.Second)
	captchaLocator := page.Locator("a.verify-bar-close")
	err = captchaLocator.WaitFor(playwright.LocatorWaitForOptions{Timeout: playwright.Float(10000)})
	if err == nil {
		fmt.Println("No captcha found.")
		captchaLocator.Click()
		textLocator := page.Locator("text=Something went wrong")
		textLocator.WaitFor(playwright.LocatorWaitForOptions{Timeout: playwright.Float(10000)})
		buttonLocator := page.Locator("text=Refresh")
		buttonLocator.WaitFor(playwright.LocatorWaitForOptions{Timeout: playwright.Float(10000)})
		buttonLocator.Click()
	}

	postLocator := page.Locator("div[data-e2e='user-post-item']")
	postCount, _ := postLocator.Count()
	fmt.Println("Total post: ", postCount)

	for i := 0; i < postCount; i++ {
		post := postLocator.Nth(i)
		post.Click()
		description, _ := post.Locator("div[data-e2e='browse-video-desc']").InnerText()
		fmt.Println("desc: ", description)
		like, _ := post.Locator("strong[data-e2e='browse-like-count']").InnerText()
		fmt.Println("like: ", like)
		comment, _ := post.Locator("strong[data-e2e='browse-comment-count']").InnerText()
		fmt.Println("comment: ", comment)
		url, _ := post.Locator("p[data-e2e='browse-video-link']").InnerText()
		fmt.Println("url: ", url)

		page.Locator("button[data-e2e='browse-close']").Click()
	}

	// targetAccount := "@kementerian.atrbpn"
	// page.Locator("form[data-e2e='search-box'] input").First().Fill(targetAccount)

	time.Sleep(10 * time.Minute)
}
