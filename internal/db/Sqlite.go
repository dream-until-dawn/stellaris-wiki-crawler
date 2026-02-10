package db

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"sync"

	_ "modernc.org/sqlite"
)

type SQLiteDB struct {
	DB *sql.DB
}

var instance *SQLiteDB
var once sync.Once

// Init initializes the SQLite database
func Init(dbPath string) {
	once.Do(func() {
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			log.Fatalf("打开sqlite数据库失败: %v", err)
		}

		// enable WAL mode for better concurrency
		_, err = db.Exec("PRAGMA journal_mode = WAL;")
		if err != nil {
			log.Printf("设置WAL模式失败: %v", err)
		}

		instance = &SQLiteDB{DB: db}

		if err := instance.migrate(); err != nil {
			log.Fatalf("数据库初始化失败: %v", err)
		}
	})
}

// Get returns the SQLite singleton instance
func Get() *SQLiteDB {
	Init("db/data.db")
	if instance == nil {
		log.Fatal("SQLite还未初始化,请先调用db.Init()")
	}
	return instance
}

const klineTableTpl = `
CREATE TABLE IF NOT EXISTS %s (
	ts     INTEGER NOT NULL,
	open   REAL    NOT NULL,
	high   REAL    NOT NULL,
	low    REAL    NOT NULL,
	close  REAL    NOT NULL,
	volume REAL    NOT NULL,
	PRIMARY KEY (ts)
) WITHOUT ROWID;
`

var tableNameRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{0,63}$`)

func isSafeTableName(name string) bool {
	return tableNameRegexp.MatchString(name)
}

func (s *SQLiteDB) migrateKline(tableName string) error {
	if !isSafeTableName(tableName) {
		return fmt.Errorf("invalid table name: %s", tableName)
	}

	stmt := fmt.Sprintf(klineTableTpl, tableName)
	_, err := s.DB.Exec(stmt)
	return err
}

func (s *SQLiteDB) migrate() error {
	tables := []string{
		"ETH_USDT_kline_1D",
		"ETH_USDT_kline_4H",
		"ETH_USDT_kline_15m",
		"BTC_USDT_kline_1D",
		"BTC_USDT_kline_4H",
		"BTC_USDT_kline_15m",
	}

	for _, t := range tables {
		if err := s.migrateKline(t); err != nil {
			return err
		}
	}
	return nil
}
