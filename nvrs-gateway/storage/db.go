package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitSQLite initializes an SQLite database connection
func InitSQLite() error {
	var err error
	DB, err = sql.Open("sqlite3", "./dev.db")
	if err != nil {
		return err
	}

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS agents (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            role TEXT NOT NULL,
			status TEXT DEFAULT 'idle'
        );

		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			agent_id INTEGER,
			task TEXT NOT NULL,
			FOREIGN KEY (agent_id) REFERENCES agents(id)
		);
    `
	if _, err = DB.Exec(createTableQuery); err != nil {
		return err
	}

	log.Println("Connected to SQLite database")
	return nil
}
