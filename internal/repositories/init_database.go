package repositories

import (
	"database/sql"
	"fmt"
	"gofinalproject/config"
	"log"
	"os"
)

func CheckAndCreateDB(dbPath string) error {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		log.Println("Database file does not exist. Creating new database.")

		db, err := sql.Open(config.DBDriver, dbPath)
		if err != nil {
			return err
		}
		defer db.Close()

		query := `
        CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date VARCHAR(8) NOT NULL, -- handle YYYYMMDD format
			title TEXT NOT NULL,
			comment TEXT DEFAULT "",
			repeat VARCHAR(128) DEFAULT ""
		);
		CREATE INDEX idx_scheduler_date ON scheduler(date);
        `
		_, err = db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
		log.Println("Database and table created successfully.")
		log.Println(dbPath, "<-- dbpath")
		return nil
	} else if err != nil {
		return err
	}

	log.Println("Database already exists.")
	log.Println(dbPath, "<-- dbpath")
	return nil
}
