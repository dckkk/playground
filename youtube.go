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

	page.Locator("ytd-rich-item-renderer").WaitFor()

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
