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
		sleepTimeSeconds = 5
	)

	// connect to the DB
	q, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")

	for {
		time.Sleep(sleepTimeSeconds * time.Second)

		runs, err := q.FetchSomeRun(context.Background())
		if err != nil {
			log.Printf("%v", err)
			continue
		}

		if len(runs) == 0 {
			continue
		}

		// we have a run!
		run := runs[0]
		log.Printf("%v", run)

		// get the config, test and solution codes
		sol, err := q.GetSolutionByID(context.Background(), run.SolutionCodeID)
		if err != nil {
			log.Printf("error getting solution code: %v for run: %v", err, run)
			continue
		}

		tes, err := q.GetTestCodeByID(context.Background(), sol.TestCodeID)
		if err != nil {
			log.Printf("error getting test code: %v for solution: %v", err, sol)
			continue
		}

		conf, err := q.GetConfByDiplayName(context.Background(), run.RunConfig)
		if err != nil {
			log.Printf("error getting conf: %v for run: %v", err, run)
			continue
		}

		// TODO: copy files to a docker container
		log.Printf("code: %v", tes.Code)
		log.Printf("solution: %v", sol.Code)
		log.Printf("Conf: %v", conf)
		// TODO: run the docker container
		// TODO: check if timeout
		// TODO: check logs / result
		// TODO: update run results
	}
}
