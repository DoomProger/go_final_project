package main

import (
	"fmt"
	"net/http"
	"time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeat := query.Get("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid 'now' date format: %v", err), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid 'date' format: %v", err), http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, date.Format("20060102"), repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
