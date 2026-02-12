package main

import (
	"log"
	"stellarisWikiCrawler/internal/crawler/browser"
	"stellarisWikiCrawler/internal/parser"

	"github.com/playwright-community/playwright-go"
)

func main() {
	bf := browser.Get()

	page, err := bf.NewPage()
	if err != nil {
		panic(err)
	}

	targetUrl := []string{
		"https://stellaris.paradoxwikis.com/Physics_research",
		"https://stellaris.paradoxwikis.com/Society_research",
		"https://stellaris.paradoxwikis.com/Engineering_research",
	}
	for i, n := 0, len(targetUrl); i < n; i++ { // 避免多次调用 length 函数。
		_, err = page.Goto(
			targetUrl[i],
			playwright.PageGotoOptions{
				WaitUntil: playwright.WaitUntilStateNetworkidle,
			},
		)
		if err != nil {
			panic(err)
		}

		result, err := parser.ParseH2PTable(page)
		if err != nil {
			panic(err)
		}

		log.Printf("url %s result: %d", targetUrl[i], len(result))
	}

	c := make(chan struct{})
	<-c
}
