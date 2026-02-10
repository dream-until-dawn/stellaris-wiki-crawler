package model

import (
	"stellarisWikiCrawler/internal/utils"

	"github.com/gocolly/colly"
)

type CrawlerCollector struct {
	Collector *colly.Collector
	Logger    *utils.Logger
}
