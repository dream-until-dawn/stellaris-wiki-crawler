package dao

func (d *DBDAO) GetTargetTree(name string) ([]TechnologyTreeItem, error) {
	rows, err := d.db.Query(`
		SELECT item.*,-CAST(closure.depth AS INTEGER) AS depth
		FROM technology_item item
		INNER JOIN technology_closure closure
		ON item.name = closure.ancestor_name
		WHERE closure.descendant_name = ?

		UNION

		SELECT item.*,closure.depth
		FROM technology_item item
		INNER JOIN technology_closure closure
		ON item.name = closure.descendant_name 
		WHERE closure.ancestor_name = ?
		ORDER BY depth ASC
	`, name, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []TechnologyTreeItem
	for rows.Next() {
		var item TechnologyTreeItem

		err := rows.Scan(
			&item.Name,
			&item.Classify,
			&item.Technology,
			&item.Icon,
			&item.Description,
			&item.Tier,
			&item.Cost,
			&item.EffectsUnlocks,
			&item.Prerequisites,
			&item.DrawWeight,
			&item.Empire,
			&item.DLC,
			&item.Depth,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	return list, nil
}

// 取出指定科技大类的直接子集关系
func (d *DBDAO) GetGraphByClassify(classify string) ([]GraphItem, error) {
	rows, err := d.db.Query(`
		SELECT 
			t.*,
			d.parent_name AS source,
			d.child_name AS target
		FROM technology_dependency d
		JOIN technology_item t
		ON t.name = d.parent_name
		WHERE t.classify = ?

		UNION ALL

		SELECT 
			t.*,
			d.parent_name AS source,
			d.child_name AS target
		FROM technology_dependency d
		JOIN technology_item t
		ON t.name = d.child_name
		WHERE t.classify = ?;
	`, classify, classify)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []GraphItem
	for rows.Next() {
		var item GraphItem

		err := rows.Scan(
			&item.Name,
			&item.Classify,
			&item.Technology,
			&item.Icon,
			&item.Description,
			&item.Tier,
			&item.Cost,
			&item.EffectsUnlocks,
			&item.Prerequisites,
			&item.DrawWeight,
			&item.Empire,
			&item.DLC,
			&item.Source,
			&item.Target,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}
	return list, nil
}

// 取出指定科技子类的直接子集关系
func (d *DBDAO) GetGraphByTechnology(technology string) ([]GraphItem, error) {
	rows, err := d.db.Query(`
		SELECT 
			t.*,
			d.parent_name AS source,
			d.child_name AS target
		FROM technology_dependency d
		JOIN technology_item t
		ON t.name = d.parent_name
		WHERE t.technology = ?

		UNION ALL

		SELECT 
			t.*,
			d.parent_name AS source,
			d.child_name AS target
		FROM technology_dependency d
		JOIN technology_item t
		ON t.name = d.child_name
		WHERE t.technology = ?;
	`, technology, technology)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []GraphItem
	for rows.Next() {
		var item GraphItem

		err := rows.Scan(
			&item.Name,
			&item.Classify,
			&item.Technology,
			&item.Icon,
			&item.Description,
			&item.Tier,
			&item.Cost,
			&item.EffectsUnlocks,
			&item.Prerequisites,
			&item.DrawWeight,
			&item.Empire,
			&item.DLC,
			&item.Source,
			&item.Target,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}
	return list, nil
}
