package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

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
