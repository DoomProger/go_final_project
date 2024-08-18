package tests

import (
	"fmt"
	"os"
	"strconv"
)

const (
	port   = 7540
	dbFile = "../scheduler.db"
)

var Port = GetPort("TODO_PORT")
var DBFile = GetDBFile("TODO_DBFILE")
var FullNextDate = false
var Search = false
var Token = ``

//TODO:
// function to file utils.go or pkg utils/utils.go

// GetPort retrieves the port number from the environment variable or uses the default value 7540.
func GetPort(envKey string) int {
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

// GetDBFile retrieves the path to db file from the environment variable or uses the default value "../scheduler.db".
func GetDBFile(envKey string) string {
	if val, ok := os.LookupEnv(envKey); ok {
		return val
	}
	return dbFile
}
