package dao

import (
	"context"
	"fmt"
)

// 返回 “全表数据 → []KLine
// asc: true 升序 false 降序
func (d *DBDAO) GetAllData(ctx context.Context, asc bool) ([]KLine, error) {
	order := "DESC"
	if asc {
		order = "ASC"
	}
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %s ORDER BY ts %s;", d.TableName, order))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]KLine, 0)
	for rows.Next() {
		f := &KLine{}
		err = rows.Scan(
			&f.Ts,
			&f.Open,
			&f.High,
			&f.Low,
			&f.Close,
			&f.Volume,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, *f)
	}
	return result, nil
}

// 返回全表最新的时间戳
func (d *DBDAO) GetLatestTs(ctx context.Context) int64 {
	row := d.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COALESCE(MAX(ts), 0) FROM %s;", d.TableName))
	var ts int64
	if err := row.Scan(&ts); err != nil {
		return 0
	}
	return ts
}

// 检查表数据的时间间隔是否正确
func (d *DBDAO) selfInspection() error {
	result, err := d.GetAllData(context.Background(), true)
	if err != nil {
		return err
	}
	timestampInterval := int64(0)
	switch d.Cycle {
	case "15m":
		timestampInterval = 1000 * 60 * 15
	case "4H":
		timestampInterval = 1000 * 60 * 60 * 4
	case "1D":
		timestampInterval = 1000 * 60 * 60 * 24
	}
	for i := 1; i < len(result); i++ {
		if result[i].Ts-result[i-1].Ts != timestampInterval {
			return fmt.Errorf("时间周期并非递增")
		}
	}
	return nil
}
