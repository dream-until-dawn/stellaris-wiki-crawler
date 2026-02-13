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
