package dao

import (
	"context"
	"fmt"
	"sort"
)

// 批量插入-主键冲突则忽视
func (d *DBDAO) BatchInsertIgnore(ctx context.Context, data []KLine) error {
	if len(data) == 0 {
		return fmt.Errorf("待插入数据为空")
	}

	// 1. 按 ts 排序，减少 B-Tree 分裂
	sort.Slice(data, func(i, j int) bool {
		return data[i].Ts < data[j].Ts
	})

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sql := fmt.Sprintf(
		`INSERT OR IGNORE INTO %s 
		(ts, open, high, low, close, volume) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		d.TableName,
	)

	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range data {
		if _, err := stmt.ExecContext(
			ctx,
			v.Ts,
			v.Open,
			v.High,
			v.Low,
			v.Close,
			v.Volume,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}
