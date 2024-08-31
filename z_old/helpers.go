package main

import (
	"fmt"
	"os"
	"strconv"
)

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

// GetDBFile retrieves the path to db file from the environment variable or uses the default value "scheduler.db".
func GetDBFile(envKey string) string {
	if val, ok := os.LookupEnv(envKey); ok {
		return val
	}
	return dbFile
}
