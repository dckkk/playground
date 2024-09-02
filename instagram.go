package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	stealth "github.com/jonfriesen/playwright-go-stealth"
	"github.com/playwright-community/playwright-go"
)

type instagramConfig struct {
	OutputCSV     string
	TargetAccount string
}

func NewInstagramConfig(
	outputCSV string,
	targetAccount string,
) instagramConfig {
	return instagramConfig{
		OutputCSV:     outputCSV,
		TargetAccount: targetAccount,
	}
}

func InstagramPOC(ctx playwright.BrowserContext, config instagramConfig) {
	minDate := time.Date(2024, 8, 26, 0, 0, 0, 0, time.UTC)
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
	targetURI := fmt.Sprintf("https://www.instagram.com/%v/", config.TargetAccount)

	if _, err := page.Goto(targetURI); err != nil {
		log.Fatalf("failed to open page %v : %v", targetURI, err)
	}

	if err := stealth.Inject(page); err != nil {
		log.Fatal("failed to inject stealth plugin")
	}

	// helper function
	parseTotal := func(v string) int {
		regex := regexp.MustCompile("\\D")
		str := regex.ReplaceAllString(v, "")
		res, _ := strconv.Atoi(str)
		return res
	}

	// Write CSV data to a file
	file, err := os.OpenFile(config.OutputCSV, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer file.Close()

	csvData := "sep=,\nusername,post,date,url,like,comment,is_repost\n"
	_, err = file.WriteString(csvData)
	if err != nil {
		log.Fatalf("could not write to file: %v", err)
	}

	page.SetDefaultTimeout(10000)
	for i := 0; i < 5; i++ {
		page.Keyboard().Down("End")
		page.Keyboard().Up("End")
		time.Sleep(time.Second * 3)
	}

	page.Keyboard().Down("Home")
	page.Keyboard().Up("Home")
	page.Locator("main > div > div:nth-child(3) > div > div").WaitFor()

	rows, err := page.Locator("main > div > div:nth-child(3) > div > div").All()
	if err != nil {
		log.Fatalf("failed to get content rows")
	}

	fetchedData := 0
	results := []Result{}
	for _, row := range rows {
		columns, err := row.Locator("> div").All()
		if err != nil {
			log.Fatalf("failed to get columns row")
			return
		}
		for _, column := range columns {
			fetchedData++
			column.Hover()
			newEntry := Result{
				Link:         "",
				TotalComment: 0,
				TotalLike:    0,
				Username:     "",
				Date:         time.Time{},
				Post:         "",
			}

			pinned := false
			if pinned, err = column.Locator("svg[aria-label=\"Pinned post icon\"]").First().IsVisible(); err != nil {
				log.Fatalf("failed to check pinned post: %v", err)
			}

			commentTotal := 0
			likeTotal := 0
			onHoverContent := []string{""}
			for i := 0; i < 5; i++ {
				// try 5 times until we can get comment and all inner texts
				onHoverContent, err = column.Locator("a").AllInnerTexts()
				if err != nil {
					log.Fatalf("failed to get like & comment: %v", err)
				}
				if len(onHoverContent[0]) > 0 {
					break
				}
			}
			fmt.Printf("\n\nonHoverContent (%v): %v\n", len(onHoverContent), onHoverContent)
			if len(onHoverContent[0]) > 0 {
				contentSplit := strings.Split(onHoverContent[0], "\n")
				if len(contentSplit) == 2 {
					likeTotal = parseTotal(contentSplit[0])
					commentTotal = parseTotal(contentSplit[1])
					fmt.Printf("comment total: %v, like total: %v \n", commentTotal, likeTotal)
					newEntry.TotalComment = commentTotal
					newEntry.TotalLike = likeTotal
				}
			}

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
			if dateTime.Before(minDate) && !pinned {
				log.Printf("already passed date: %v", dateTime)
				fmt.Printf("post dumped to Csv: %v", fetchedData-1)
				return
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
