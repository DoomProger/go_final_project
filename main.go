package main

import (
	"fmt"
	"gofinalproject/tests"
	"log"
	"net/http"
)

// var Port = tests.GetPort("TODO_PORT", 7540)

func main() {
	Port := tests.GetPort("TODO_PORT", 7540)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
