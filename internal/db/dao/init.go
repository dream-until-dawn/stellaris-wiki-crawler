package dao

import (
	"database/sql"

	"stellarisWikiCrawler/internal/db"
)

type KLine struct {
	Ts     int64   `json:"ts"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

type DBDAO struct {
	InstId    string
	Cycle     string
	TableName string
	db        *sql.DB
}

func NewDBDAO(instId string, cycle string) *DBDAO {
	tableName := instId + "_USDT_kline_" + cycle
	dbDao := &DBDAO{
		InstId:    instId,
		Cycle:     cycle,
		TableName: tableName,
		db:        db.Get().DB,
	}
	err := dbDao.selfInspection()
	if err != nil {
		panic(err)
	}
	return dbDao
}
