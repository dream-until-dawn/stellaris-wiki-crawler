package main

import (
	"log"
	"stellarisWikiCrawler/internal/crawler"
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

	_, err = page.Goto(
		"https://stellaris.paradoxwikis.com/Society_research",
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

	log.Printf("result: %d", len(result))
	c := make(chan struct{})
	<-c
}

func main_1() {
	c := crawler.NewCollector()

	// 科技介绍
	// c.Collector.Visit("https://stellaris.paradoxwikis.com/Technology")
	// 物理学
	c.Collector.Visit("https://stellaris.paradoxwikis.com/Physics_research")
	// 社会学
	// c.Collector.Visit("https://stellaris.paradoxwikis.com/Society_research")
	// 工程学
	// c.Collector.Visit("https://stellaris.paradoxwikis.com/Engineering_research")

	c.Collector.Wait()
}
