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
	ryears  int
	rmonths int
	rdays   int
}

func checkAndCreateDB(dbPath string) error {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		log.Println("Database file does not exist. Creating new database.")
		db, err := sql.Open("sqlite3", dbPath)
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
		// fmt.Printf("Unsupported value %v", repeat)
		return repeatDate{}, fmt.Errorf("input value is empty [%s]", repeat)
	}

	repeatSettings := strings.Split(repeat, " ")
	if repeatSettings[0] == "y" && len(repeatSettings) == 1 {
		return repeatDate{
			ryears: 1,
		}, nil
	}

	if repeatSettings[0] == "d" && len(repeatSettings) == 2 {
		v, err := strconv.Atoi(repeatSettings[1])
		if v == 0 {
			return repeatDate{}, fmt.Errorf("the number of repeating days must be greater than 0")
		}
		if err != nil {
			return repeatDate{}, fmt.Errorf("the number of repeating days must be define")
		}

		if v > 400 {
			return repeatDate{}, fmt.Errorf("maximum number of repeating days is 400 but got %d", v)
		}

		return repeatDate{
			rdays: v,
		}, nil
	}
	return repeatDate{}, fmt.Errorf("can't parse repeat values, got %v", repeatSettings)
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	d, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("failed parse date %w", err)
	}

	rdate, err := getRepeat(repeat)
	if err != nil {
		return "", err
	}

	nextDate := d.AddDate(rdate.ryears, rdate.rmonths, rdate.rdays)

	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(rdate.ryears, rdate.rmonths, rdate.rdays)
	}

	resNextDate := nextDate.Format("20060102")
	// fmt.Println(resNextDate)
	return resNextDate, nil
}
