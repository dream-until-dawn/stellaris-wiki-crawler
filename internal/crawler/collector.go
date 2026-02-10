package crawler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocolly/colly"

	"stellarisWikiCrawler/internal/utils"
)

type CrawlerCollector struct {
	Collector *colly.Collector
	logger    *utils.Logger
}

// ================================
// NewCollector
// ================================
func NewCollector() *CrawlerCollector {
	c := colly.NewCollector(
		// Stellaris Wiki 域名
		colly.AllowedDomains(
			"stellaris.paradoxwikis.com",
			"www.stellaris.paradoxwikis.com",
		),

		// 异步爬取（非常重要）
		colly.Async(true),

		// 自动去重（基于 URL）
		colly.CacheDir(".colly_cache"),
	)

	// ----------------------------
	// 基础配置
	// ----------------------------
	c.WithTransport(&http.Transport{
		DisableKeepAlives: false,
		MaxIdleConns:      50,
		IdleConnTimeout:   30 * time.Second,
	})

	c.UserAgent = "StellarisWikiCrawler/1.0 (+https://github.com/yourname/stellaris-wiki)"
	// ----------------------------
	// logger
	// ----------------------------
	l := utils.NewLogger().WithTask("爬虫")
	cc := &CrawlerCollector{
		Collector: c,
		logger:    l,
	}

	// ----------------------------
	// 限速器（来自 limiter.go）
	// ----------------------------
	ApplyRateLimiter(c, l)

	// ----------------------------
	// 全局回调
	// ----------------------------
	cc.registerCallbacks()

	return cc
}

// ================================
// Register Callbacks
// ================================
func (c *CrawlerCollector) registerCallbacks() {

	// 请求前
	c.Collector.OnRequest(func(r *colly.Request) {
		c.logger.Info(fmt.Sprintf("请求URL: %s", r.URL.String()))
	})

	// 响应成功
	c.Collector.OnResponse(func(r *colly.Response) {
		c.logger.Suc(
			fmt.Sprintf(
				"请求成功: %d %s (%d bytes)",
				r.StatusCode, r.Request.URL.String(), len(r.Body),
			),
		)
		log.Printf("==== HTML BEGIN (%s) ====", r.Request.URL)
		log.Println(string(r.Body))
		log.Printf("==== HTML END ====")
	})

	// HTML 解析入口（统一路由）
	c.Collector.OnHTML("html", func(e *colly.HTMLElement) {
		c.Route(e)
	})

	// 请求错误
	c.Collector.OnError(func(r *colly.Response, err error) {
		if r != nil {
			c.logger.Error(
				fmt.Sprintf(
					"请求失败: %s status=%d err=%v",
					r.Request.URL.String(), r.StatusCode, err,
				),
			)
		} else {
			c.logger.Error(
				fmt.Sprintf(
					"请求失败: err=%v", err,
				),
			)
		}
	})

	// 请求完成
	c.Collector.OnScraped(func(r *colly.Response) {
		c.logger.Debug(
			fmt.Sprintf(
				"请求完成: %s", r.Request.URL.String(),
			),
		)
	})
}
