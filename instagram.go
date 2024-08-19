package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

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
