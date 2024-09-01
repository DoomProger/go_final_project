package config

import (
	"fmt"
	"os"
	"strconv"
)

// DB and http server default settings
const (
	DateFormat       = "20060102" //YYYYMMDD
	DateFormatSearch = "02.01.2006"
	DBDriver         = "sqlite3"
	dbFile           = "scheduler.db"
	port             = 7540
	todoPassword     = "123"
)

// SQL limit query const
const (
	Limit50 int = 50
)

// Login settings
const (
	TokenTTL = 8
)

var TodoPassword = getPassword("TODO_PASSWORD")

var DBFile = getDBFile("TODO_DBFILE")
var Port = getPort("TODO_PORT")

// GetPort retrieves the port number from the environment variable or uses the default value 7540.
func getPort(envKey string) int {
	if val, ok := os.LookupEnv(envKey); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return port
		}
		return v
	}
	return port
}

// GetDBFile retrieves the path to db file from the environment variable or uses the default value "scheduler.db".
func getDBFile(envKey string) string {
	if val, ok := os.LookupEnv(envKey); ok {
		return val
	}
	return dbFile
}

func getPassword(envKey string) string {
	if val, ok := os.LookupEnv(envKey); ok {
		return val
	}
	return todoPassword
}
