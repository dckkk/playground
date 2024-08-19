package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
)

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
