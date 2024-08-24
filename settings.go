package main

// DB and http server default settings
const (
	dateFormat       = "20060102" //YYYYMMDD
	dateFormatSearch = "02.01.2006"
	// dateFormatSearch = "01.02.2006"
	DBDriver = "sqlite3"
	dbFile   = "scheduler.db"
	port     = 7540
)

// SQL limit query const
const (
	limit50  int = 50
	limit100 int = 100
)

// Login settings
const (
	Token = ``
)

var DBFile = GetDBFile("TODO_DBFILE")
var Port = GetPort("TODO_PORT")
