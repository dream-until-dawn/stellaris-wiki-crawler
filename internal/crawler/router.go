package crawler

import (
	"fmt"
	"net/url"
	"stellarisWikiCrawler/internal/parser"
	"strings"

	"github.com/gocolly/colly"
)

// ================================
// Route Entry
// ================================
func (c *CrawlerCollector) Route(e *colly.HTMLElement) {
	rawURL := e.Request.URL.String()

	// 过滤无意义页面
	if shouldIgnore(rawURL) {
		return
	}

	switch {
	case isTechOverviewPage(rawURL):
		c.logger.Info(
			fmt.Sprintf("科技概述页面: %s", rawURL),
		)
		// ParseTechOverview(e)

	case isTechListPage(rawURL):
		c.logger.Info(
			fmt.Sprintf("科技列表页面: %s", rawURL),
		)
		parser.ParseTechList(e, c.logger)

	case isTechDetailPage(rawURL):
		c.logger.Info(
			fmt.Sprintf("科技详情页面: %s", rawURL),
		)
		// ParseTechDetail(e)

	default:
		// 非目标页面不处理
		c.logger.Warn(
			fmt.Sprintf("未处理的页面: %s", rawURL),
		)
	}
}

// ================================
// Ignore Rules
// ================================
func shouldIgnore(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return true
	}

	// MediaWiki 行为页
	if strings.Contains(u.Path, "/Special:") ||
		strings.Contains(u.Path, "/File:") ||
		strings.Contains(u.Path, "/Help:") ||
		strings.Contains(u.Path, "/Template:") ||
		strings.Contains(u.Path, "/Talk:") {
		return true
	}

	// 编辑、历史、diff 等
	q := u.Query()
	if q.Has("action") {
		return true
	}

	return false
}

// ================================
// Page Type Detection
// ================================

// 科技总览页（入口）
func isTechOverviewPage(rawURL string) bool {
	return strings.HasSuffix(rawURL, "/Technology")
}

// 各领域科技列表页
// Physics / Society / Engineering
func isTechListPage(rawURL string) bool {

	return strings.Contains(rawURL, "Physics") ||
		strings.Contains(rawURL, "Society") ||
		strings.Contains(rawURL, "Engineering")
}

// 单个科技详情页
func isTechDetailPage(rawURL string) bool {
	// Wiki 中科技详情一般不在 Technology 目录下
	// 例如： /Research_Stations
	if strings.Contains(rawURL, "/Technology") {
		return false
	}

	// 排除明显不是科技的页面
	blacklist := []string{
		"/Traditions",
		"/Ascension_perks",
		"/Civics",
		"/Ethics",
	}

	for _, b := range blacklist {
		if strings.Contains(rawURL, b) {
			return false
		}
	}

	return true
}
