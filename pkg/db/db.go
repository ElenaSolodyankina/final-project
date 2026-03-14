package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR NOT NULL DEFAULT '',
    comment TEXT,
    "repeat" VARCHAR(128) DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	fileExists := err == nil

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if !fileExists {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
		log.Println("[DB] Database created:", dbFile)
	} else {
		log.Println("[DB] Database accessed:", dbFile)
	}

	return nil
}
