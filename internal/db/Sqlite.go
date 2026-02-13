package db

import (
	"database/sql"
	"fmt"
	"log"
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

		// 设置WAL模式
		_, err = db.Exec("PRAGMA journal_mode = WAL;")
		if err != nil {
			log.Printf("设置WAL模式失败: %v", err)
		}

		// 启用外键
		_, err = db.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			log.Printf("启用外键失败: %v", err)
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
		log.Fatal("SQLite还未初始化")
	}
	return instance
}

func (s *SQLiteDB) creationTechnologyTable() error {
	const tableTpl = `
		CREATE TABLE IF NOT EXISTS technology (
			classify TEXT NOT NULL,
			name TEXT PRIMARY KEY,
			description TEXT NOT NULL
		);
	`
	stmt := fmt.Sprintf(tableTpl)
	_, err := s.DB.Exec(stmt)
	return err
}

func (s *SQLiteDB) createTechnologyItemTable() error {
	const createTableSQL = `
		CREATE TABLE IF NOT EXISTS technology_item (
			name TEXT PRIMARY KEY,
			classify TEXT,
			technology TEXT,
			icon TEXT,
			description TEXT,
			tier TEXT,
			cost TEXT,
			effects_unlocks TEXT,
			prerequisites TEXT,
			draw_weight TEXT,
			empire TEXT,
			dlc TEXT,
			FOREIGN KEY (technology) REFERENCES technology(name)
		) WITHOUT ROWID;
	`
	_, err := s.DB.Exec(createTableSQL)
	return err
}
func (s *SQLiteDB) createTechnologyDependencyTable() error {
	const createTableSQL = `
		CREATE TABLE IF NOT EXISTS technology_dependency (
			parent_name TEXT NOT NULL,
			child_name  TEXT NOT NULL,
			PRIMARY KEY (parent_name, child_name),
			FOREIGN KEY (parent_name) REFERENCES technology_item(name),
			FOREIGN KEY (child_name) REFERENCES technology_item(name)
		) WITHOUT ROWID;
	`
	_, err := s.DB.Exec(createTableSQL)
	return err
}
func (s *SQLiteDB) createTechnologyClosureTable() error {
	const createTableSQL = `
		CREATE TABLE IF NOT EXISTS technology_closure (
			ancestor_name   TEXT NOT NULL,
			descendant_name TEXT NOT NULL,
			depth           INTEGER NOT NULL,
			PRIMARY KEY (ancestor_name, descendant_name),
			FOREIGN KEY (ancestor_name) REFERENCES technology_item(name),
			FOREIGN KEY (descendant_name) REFERENCES technology_item(name)
		) WITHOUT ROWID;
	`
	_, err := s.DB.Exec(createTableSQL)
	return err
}

func (s *SQLiteDB) createTechnologyClosureIndexes() error {
	const sqlStmt = `
		CREATE INDEX IF NOT EXISTS idx_closure_ancestor
		ON technology_closure(ancestor_name);

		CREATE INDEX IF NOT EXISTS idx_closure_descendant
		ON technology_closure(descendant_name);

		CREATE INDEX IF NOT EXISTS idx_closure_desc_depth
		ON technology_closure(descendant_name, depth);
	`
	_, err := s.DB.Exec(sqlStmt)
	return err
}

func (s *SQLiteDB) migrate() error {
	if err := s.creationTechnologyTable(); err != nil {
		return err
	}
	if err := s.createTechnologyItemTable(); err != nil {
		return err
	}
	if err := s.createTechnologyDependencyTable(); err != nil {
		return err
	}
	if err := s.createTechnologyClosureTable(); err != nil {
		return err
	}
	if err := s.createTechnologyClosureIndexes(); err != nil {
		return err
	}
	return nil
}
