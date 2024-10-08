package main

import (
	"fmt"
	"log"

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

	// igAccount := []string{
	// 	"kanwilbpndkijakarta",
	// 	"kantahkotajakpus",
	// 	"kantahkotajakartatimur",
	// 	"kantahkotajakartabarat",
	// 	"kantahkotajakartautara",
	// 	"kantahkotajakartaselatan",
	// }

	ytAccount := []string{
		"@kanwilbpndkijakarta4650",
		"@kantahkotajakartapusat",
		"@kantahkotajakartatimur991",
		"@kantahkotajakartabarat",
		"@kantahkotajakartautara3691",
		"@kantahkotajakartaselatan",
	}

	// FacebookPOC(context)
	// TwitterPOC(context)
	// for _, v := range igAccount {
	// 	InstagramPOC(context, NewInstagramConfig(
	// 		fmt.Sprintf("output/instagram_%v.csv", v),
	// 		v,
	// 	))
	// }
	for _, v := range ytAccount {
		YoutubePOC(context, NewYoutubeConfig(
			fmt.Sprintf("output/youtube_%v.csv", v),
			v,
		))
	}
	// tiktokPOC(context)
}
