package main

import (
	"context"
	"log"
	"time"

	"../sql/db"
)

func main() {
	// configuration consts
	const (
		sleepTimeSeconds = 2
	)

	// connect to the DB
	q, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")

	for {
		time.Sleep(sleepTimeSeconds * time.Second)

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
