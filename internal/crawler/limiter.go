package crawler

import (
	"log"
	"time"

	"github.com/gocolly/colly"

	"stellarisWikiCrawler/internal/config"
	"stellarisWikiCrawler/internal/utils"
)

// ================================
// Apply Rate Limiter
// ================================
func ApplyRateLimiter(c *colly.Collector, l *utils.Logger) {
	cfg := config.Get()

	parallelism := cfg.CrawlerParallelism
	delay, err := time.ParseDuration(cfg.CrawlerDelay)
	randomDelay, err := time.ParseDuration(cfg.CrawlerRandomDelay)
	if err != nil {
		log.Fatalf("读取爬虫配置失败: %v", err)
	}

	// 应用到 Colly
	err = c.Limit(&colly.LimitRule{
		DomainGlob:  "*stellaris*wiki*",
		Parallelism: parallelism,
		Delay:       delay,
		RandomDelay: randomDelay,
	})

	if err != nil {
		log.Fatalf("设置爬虫限速失败: %v", err)
	} else {
		l.Suc("设置爬虫限速成功")
	}
}
