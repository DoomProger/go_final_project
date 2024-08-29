package nextdate

import (
	"fmt"
	"gofinalproject/config"
	"time"
)

type repeatDate struct {
	years  int
	months int
	days   int
}

func NormalizeTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	now = NormalizeTime(now)
	startDate, err := time.Parse(config.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("failed parse date %w", err)
	}

	rdate, err := getRepeat(repeat)
	if err != nil {
		return "", err
	}

	if !startDate.Before(now) {
		nextDate := startDate.AddDate(rdate.years, rdate.months, rdate.days)
		return nextDate.Format(config.DateFormat), nil
	}

	for startDate.Before(now) {
		startDate = startDate.AddDate(rdate.years, rdate.months, rdate.days)
	}

	nextDate := startDate.Format(config.DateFormat)
	return nextDate, nil
}
