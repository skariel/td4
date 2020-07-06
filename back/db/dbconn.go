// Package db handling all queries to DB
package db

import (
	"database/sql"
	"log"

	gdb "td4/back/db/generated"

	// for talking to postgres
	_ "github.com/lib/pq"
)

// ConnectDB establish the global DB connection
func ConnectDB() (*gdb.Queries, *sql.DB) {
	connStr := `user=postgres
				dbname=skariel
				password=1234567
				host=127.0.0.1
				port=5432`
	dbc, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	err = dbc.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return gdb.New(dbc), dbc
}
