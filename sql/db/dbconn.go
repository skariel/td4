// Package db handling all queries to DB
package db

import (
	"database/sql"

	// for talking to postgres
	_ "github.com/lib/pq"
)

// ConnectDB establish the global DB connection
func ConnectDB() (*Queries, error) {
	connStr := `user=postgres
				dbname=skariel
				password=1234567
				host=127.0.0.1
				port=5432`
	dbc, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = dbc.Ping()
	if err != nil {
		return nil, err
	}

	return New(dbc), nil
}
