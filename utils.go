package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type repeatDate struct {
	years  int
	months int
	days   int
}

func checkAndCreateDB(dbPath string) error {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		log.Println("Database file does not exist. Creating new database.")
		db, err := sql.Open(DBDriver, dbPath)
		if err != nil {
			return err
		}
		defer db.Close()

		query := `
        CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date VARCHAR(8) DEFAULT "", -- handle YYYYMMDD format
			title TEXT DEFAULT "",
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

func getRepeat(repeat string) (repeatDate, error) {
	if len(repeat) == 0 {
		return repeatDate{}, fmt.Errorf("input value is empty [%s]", repeat)
	}

	repeatSettings := strings.Split(repeat, " ")

	if repeatSettings[0] == "y" && len(repeatSettings) == 1 {
		return repeatDate{
			years: 1,
		}, nil
	}

	if repeatSettings[0] == "d" && len(repeatSettings) == 2 {
		v, err := strconv.Atoi(repeatSettings[1])
		if err != nil {
			return repeatDate{}, fmt.Errorf("the number of repeating days must be define")
		}

		if v == 0 || v > 400 {
			return repeatDate{}, fmt.Errorf("the number of repeating days must be between 0 and 400, but go [%d]", v)
		}

		return repeatDate{
			days: v,
		}, nil
	}
	return repeatDate{}, fmt.Errorf("can't parse repeat values, got %v", repeatSettings)
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	startDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", fmt.Errorf("failed parse date %w", err)
	}

	rdate, err := getRepeat(repeat)
	if err != nil {
		return "", err
	}

	if !startDate.Before(now) && rdate.years == 1 {
		nextDate := startDate.AddDate(rdate.years, rdate.months, rdate.days)
		return nextDate.Format(dateFormat), nil
	}

	if !startDate.Before(now) && rdate.days == 1 {
		nextDate := startDate.AddDate(rdate.years, rdate.months, rdate.days)
		return nextDate.Format(dateFormat), nil
	}

	if !startDate.Before(now) && rdate.days > 1 {
		nextDate := startDate.AddDate(rdate.years, rdate.months, rdate.days)
		return nextDate.Format(dateFormat), nil
	}

	for startDate.Before(now) {
		startDate = startDate.AddDate(rdate.years, rdate.months, rdate.days)
	}

	nextDate := startDate.Format(dateFormat)
	return nextDate, nil
}
