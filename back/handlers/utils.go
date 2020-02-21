package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"../../sql/db"
)

var (
	q *db.Queries
)

// SetQueries globally set the package query object from sqlc
func SetQueries(_q *db.Queries) {
	q = _q
}

func ise(w http.ResponseWriter, err error) {
	log.Printf("%v", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("something bad happened!"))
}

func forb(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Forbidden"))
}

func rj(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		ise(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func rh(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(s))
}
