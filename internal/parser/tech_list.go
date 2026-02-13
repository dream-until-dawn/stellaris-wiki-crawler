package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"stellarisWikiCrawler/internal/model"
	"stellarisWikiCrawler/internal/scripts"
	"strings"

	"github.com/playwright-community/playwright-go"
)

// 解析页面
func ParseTechList(page playwright.Page) ([]model.CrawierTechnology, error) {
	result, err := page.Evaluate(scripts.ParseH2PTableJs)
	if err != nil {
		return nil, err
	}
	// 先把 interface{} 转成 json
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	// 再反序列化到你的结构体
	var tech []model.CrawierTechnology
	if err := json.Unmarshal(data, &tech); err != nil {
		return nil, err
	}
	// 判断 error 字段
	if len(tech) == 0 {
		return nil, fmt.Errorf("序列化失败")
	}

	// 保存为json方便查看
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}
	url := page.URL()
	idx := strings.LastIndex(url, "/")
	save_err := os.WriteFile(fmt.Sprintf("env/%s.json", url[idx+1:]), b, 0644)
	if save_err != nil {
		return nil, save_err
	}
	return tech, nil
}
