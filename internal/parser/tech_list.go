package parser

import (
	"fmt"
	"stellarisWikiCrawler/internal/utils"

	"github.com/gocolly/colly"
)

func ParseTechList(e *colly.HTMLElement, l *utils.Logger) {
	l.Debug("开始处理页面内容")

	e.ForEach("a[href]", func(_ int, a *colly.HTMLElement) {
		link := a.Attr("href")
		fmt.Println("发现链接:", link)
	})
}
