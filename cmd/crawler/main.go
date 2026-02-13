package main

import (
	"log"
	"stellarisWikiCrawler/internal/crawler/browser"
	"stellarisWikiCrawler/internal/db/dao"
	"stellarisWikiCrawler/internal/parser"

	"github.com/playwright-community/playwright-go"
)

type TargetUrl struct {
	Url      string
	Classify string
}

func main() {
	dao := dao.NewDBDAO()
	bf := browser.Get()

	page, err := bf.NewPage()
	if err != nil {
		panic(err)
	}

	targetUrl := []TargetUrl{
		{
			Url:      "https://stellaris.paradoxwikis.com/Physics_research",
			Classify: "physics",
		},
		{
			Url:      "https://stellaris.paradoxwikis.com/Society_research",
			Classify: "society",
		},
		{
			Url:      "https://stellaris.paradoxwikis.com/Engineering_research",
			Classify: "engineering",
		},
	}
	for i, n := 0, len(targetUrl); i < n; i++ {
		_, err = page.Goto(
			targetUrl[i].Url,
			playwright.PageGotoOptions{
				WaitUntil: playwright.WaitUntilStateNetworkidle,
			},
		)
		if err != nil {
			panic(err)
		}

		result, err := parser.ParseTechList(page)
		if err != nil {
			panic(err)
		}
		amount := 0
		for _, t := range result {
			amount += len(t.Table)
		}

		log.Printf("url %s 结果对象长度: %d ,数目: %d\n", targetUrl[i], len(result), amount)
		dbErr := dao.BuildTechTree(result, targetUrl[i].Classify)
		if dbErr != nil {
			log.Printf("数据库操作失败: %v\n", dbErr)
		}
	}
}
