package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/livingpool/constants"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DB *sql.DB
}

func Setup(isProd bool) (*Config, error) {
	if isProd {
		url := os.Getenv("TURSO_DATABASE_URL") + "?authToken=" + os.Getenv("TURSO_AUTH_TOKEN")
		db, err := sql.Open("libsql", url)
		if err != nil {
			return nil, fmt.Errorf("error connecting to turso: %v", err)
		}

		var sqlStmt string
		for i := constants.DIGIT_LOWER_LIMIT; i <= constants.DIGIT_UPPER_LIMIT; i++ {
			stmt := fmt.Sprintf("create table if not exists board%d (id integer primary key, name text unique, attempts integer);", i)
			sqlStmt += stmt
		}

		if _, err = db.Exec(sqlStmt); err != nil {
			return nil, fmt.Errorf("error executing sql: %v", err)
		}

		return &Config{DB: db}, nil
	} else {
		db, err := sql.Open("sqlite3", "file:test.db?mode=memory&cache=shared")
		if err != nil {
			return nil, err
		}

		if err := populateLocalDB(db); err != nil {
			return nil, err
		}

		return &Config{DB: db}, nil
	}
}

func populateLocalDB(db *sql.DB) error {
	for i := constants.DIGIT_LOWER_LIMIT; i <= constants.DIGIT_UPPER_LIMIT; i++ {
		if _, err := db.Exec(fmt.Sprintf("create table board%d (id integer primary key, name text unique, attempts integer)", i)); err != nil {
			return fmt.Errorf("create table board%d failed: %v", i, err)
		}
	}

	for i := constants.DIGIT_LOWER_LIMIT; i <= constants.DIGIT_UPPER_LIMIT; i++ {
		for j := range constants.MAX_ROWS_DISPLAYED + 1 {
			if _, err := db.Exec(fmt.Sprintf("insert into board%d (name, attempts) values ('tim%[2]d', %[2]d)", i, j+10)); err != nil {
				return fmt.Errorf("insert into board%d failed: %v", i, err)
			}
		}
	}

	return nil
}
