package dao

import (
	"database/sql"
	"stellarisWikiCrawler/internal/model"
	"stellarisWikiCrawler/internal/utils"
)

func (d *DBDAO) InsertTechnology(item model.CrawierTechnology, classify string) error {
	// 1️⃣ 插入 technology 表
	_, err := d.db.Exec(`
		INSERT INTO technology(name, classify, description)
		VALUES (?, ?, ?)
		ON CONFLICT(name)
		DO UPDATE SET
			classify = excluded.classify,
			description = excluded.description
	`, item.Title, classify, item.Description)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBDAO) InsertTechnologyItem(technology string, classify string, item model.CrawierTechnologyItem) error {
	name := item.Technology.Name
	// 1️⃣ 插入 technology_item 表
	_, err := d.db.Exec(`
		INSERT INTO technology_item(name, classify, technology, icon, description, tier, cost, effects_unlocks, prerequisites, draw_weight, empire, dlc)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(name)
		DO UPDATE SET
			classify = excluded.classify,
			technology = excluded.technology,
			icon = excluded.icon,
			description = excluded.description,
			tier = excluded.tier,
			cost = excluded.cost,
			effects_unlocks = excluded.effects_unlocks,
			prerequisites = excluded.prerequisites,
			draw_weight = excluded.draw_weight,
			empire = excluded.empire,
			dlc = excluded.dlc
	`,
		name,
		classify,
		technology,
		item.Icon,
		item.Technology.Description,
		item.Tier,
		item.Cost,
		item.EffectsUnlocks,
		item.Prerequisites,
		item.DrawWeight,
		item.Empire,
		item.DLC,
	)
	if err != nil {
		return err
	}
	// 2️⃣ 插入自己 -> 自己
	_, err = d.db.Exec(`
		INSERT OR IGNORE INTO technology_closure
		(ancestor_name, descendant_name, depth)
		VALUES (?, ?, 0)
	`, name, name)
	if err != nil {
		return err
	}
	// 3️⃣ 处理父节点
	for _, parent := range item.Prerequisites {
		prerequisites := utils.CleanTechName(parent)
		// 检查父节点是否存在
		var exists string
		err = d.db.QueryRow(
			`SELECT name FROM technology_item WHERE name = ?`,
			prerequisites,
		).Scan(&exists)

		if err == sql.ErrNoRows {
			continue // 父节点不存在，跳过
		}
		if err != nil {
			return err
		}

		// 4️⃣ 插入直接依赖
		_, err = d.db.Exec(`
			INSERT OR IGNORE INTO technology_dependency(parent_name, child_name)
			VALUES (?, ?)
		`, prerequisites, name)
		if err != nil {
			return err
		}

		// 5️⃣ 复制父节点的所有祖先路径
		_, err = d.db.Exec(`
			INSERT OR IGNORE INTO technology_closure
			(ancestor_name, descendant_name, depth)
			SELECT ancestor_name, ?, depth + 1
			FROM technology_closure
			WHERE descendant_name = ?
		`, name, prerequisites)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DBDAO) BuildTechTree(data []model.CrawierTechnology, classify string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, tech := range data {
		err = d.InsertTechnology(tech, classify)
		if err != nil {
			return err
		}
		for _, item := range tech.Table {
			err = d.InsertTechnologyItem(tech.Title, classify, item)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
