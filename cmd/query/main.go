package main

import (
	"fmt"
	"stellarisWikiCrawler/internal/db/dao"
)

func main() {
	dao := dao.NewDBDAO()

	tree, err := dao.GetTargetTree("Weaver Organo-Tech Amplifier")
	if err != nil {
		fmt.Printf("查询出错: %+v", err)
	}
	fmt.Printf("查询结果: %+v", tree)
}
