package tests

import (
	"fmt"
	"os"
	"strconv"
)

// var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = false
var Token = ``

// GetPort retrieves the port number from the environment variable or uses the default value.
func GetPort(envKey string, def int) int {
	if val, ok := os.LookupEnv(envKey); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return def
		}
		return v
	}
	return def
}

// var envKey = "TODO_PORT"
// var Port = getEnv("TODO_PORT", 7540)
