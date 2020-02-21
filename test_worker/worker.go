package main

import (
	"context"
	"log"
	"time"

	"../sql/db"
)

var q *db.Queries

func main() {
	// connect to the DB
	q, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB")

	for true {
		time.Sleep(2 * time.Second)
		run, err := q.FetchSomeRun(context.Background())
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		if len(run) > 0 {
			log.Printf("%v", run)
		}
	}
}
