package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"stellarisWikiCrawler/internal/scripts"
	"stellarisWikiCrawler/internal/utils"

	"github.com/gocolly/colly"
	"github.com/playwright-community/playwright-go"
)

func ParseTechList(e *colly.HTMLElement, l *utils.Logger) {
	l.Debug("开始处理页面内容")
}

// 解析起是否是h2 + p + table结构 (且其概率嵌套 h3 + p + table)
func ParseH2PTable(page playwright.Page) ([]map[string]string, error) {
	_, err := page.Evaluate(scripts.ParseH2PTableJs)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ParseTableByIndex
// page: Playwright 页面对象
// tableIndex: 第几个 table (从 0 开始)
// 返回：[]map[string]string
func ParseTableByIndex(page playwright.Page, tableIndex int) ([]map[string]string, error) {
	result, err := page.Evaluate(`
		(tableIndex) => {
			// 获取全部 table 标签
			const tables = document.querySelectorAll("table")
			if (!tables || tables.length <= tableIndex) {
				return { error: "table not found" }
			}
			console.log('tables',tables);
			const table = tables[tableIndex]
			// 获取表头
			let headers = []
			const thead = table.querySelector("thead")
			if (thead) {
				headers = Array.from(thead.querySelectorAll("th"))
					.map(th => th.innerText.trim())
			}
			// 获取 tbody
			const tbody = table.querySelector("tbody")
			if (!tbody) {
				return []
			}
			const rows = Array.from(tbody.querySelectorAll("tr"))
			return rows.map(row => {
				const cells = Array.from(row.querySelectorAll("td"))
				let obj = {}
				cells.forEach((cell, index) => {
					let key = ""
					if (headers.length > 0) {
						key = headers[index] || ("col_" + index)
					} else {
						key = "col_" + index
					}
					obj[key] = cell.innerHTML.trim()
				})
				return obj
			})
		}
	`, tableIndex)

	if err != nil {
		return nil, err
	}

	// JS返回错误
	if m, ok := result.(map[string]interface{}); ok {
		if errMsg, exists := m["error"]; exists {
			return nil, fmt.Errorf("%v", errMsg)
		}
	}

	rows := result.([]interface{})
	var final []map[string]string

	for _, r := range rows {
		row := r.(map[string]interface{})
		m := make(map[string]string)
		for k, v := range row {
			if str, ok := v.(string); ok {
				m[k] = str
			} else {
				m[k] = fmt.Sprintf("%v", v)
			}
		}

		final = append(final, m)
	}

	// 保存为json方便查看
	b, err := json.MarshalIndent(final, "", "  ")
	if err != nil {
		return nil, err
	}
	save_err := os.WriteFile("env/tech_list.json", b, 0644)
	if save_err != nil {
		return nil, save_err
	}

	return final, nil
}
