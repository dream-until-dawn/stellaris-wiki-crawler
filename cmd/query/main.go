package main

import (
	"encoding/json"
	"fmt"
	"os"
	"stellarisWikiCrawler/internal/db/dao"
)

func main() {
	dao := dao.NewDBDAO()

	// res, err := dao.GetTargetTree("Weaver Organo-Tech Amplifier")
	// res, err := dao.GetGraphByClassify("physics")
	res, err := dao.GetGraphByTechnology("Biology")
	if err != nil {
		fmt.Printf("查询出错: %+v", err)
	}
	// fmt.Printf("查询结果: %+v", res)

	// 保存为json方便查看
	b, err := json.MarshalIndent(res, "", "  ")
	os.WriteFile(fmt.Sprintf("vue/public/%s.json", "test_links"), b, 0644)
}
