package main

import (
	"os"
)

var (
	GitCommit = ""
	GitBranch = "dirty"
)

func main() {
	if err := InitDB(os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PORT")); err != nil {
		panic(err)
	}
	StartServer(os.Getenv("PORT"))
}
