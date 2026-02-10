package browser

import (
	"log"
	"sync"

	"stellarisWikiCrawler/internal/utils"

	"github.com/playwright-community/playwright-go"
)

type BrowserFetcher struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	logger  *utils.Logger
}

var (
	instance *BrowserFetcher
	once     sync.Once
)

// Get 返回全局 BrowserFetcher 单例
func Get() *BrowserFetcher {
	once.Do(func() {
		instance = newBrowserFetcher()
	})
	return instance
}

func newBrowserFetcher() *BrowserFetcher {
	logger := utils.NewLogger().WithTask("无头浏览器")
	// 启动 Playwright
	if err := playwright.Install(); err != nil {
		log.Fatalf("playwright 安装失败: %v", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("playwright 运行失败: %v", err)
	}

	// 启动 Chromium（默认就能过大多数 JS Challenge）
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // 调试时可改成 false
		Args: []string{
			"--disable-blink-features=AutomationControlled",
			"--no-sandbox",
		},
	})
	if err != nil {
		log.Fatalf("browser 启动失败: %v", err)
	}

	logger.Suc("browser 启动")

	return &BrowserFetcher{
		pw:      pw,
		browser: browser,
		logger:  logger,
	}
}

// NewPage 创建一个新的页面（带常用反爬配置）
func (b *BrowserFetcher) NewPage() (playwright.Page, error) {
	context, err := b.browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String(
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
				"(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		),
		Locale: playwright.String("en-US"),
	})
	if err != nil {
		return nil, err
	}

	page, err := context.NewPage()
	if err != nil {
		return nil, err
	}

	return page, nil
}

// Close 用于程序退出时调用
func (b *BrowserFetcher) Close() {
	if b.browser != nil {
		_ = b.browser.Close()
	}
	if b.pw != nil {
		_ = b.pw.Stop()
	}
	b.logger.Info("browser 停止")
}
