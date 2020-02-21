package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"../../sql/db"
)

type key int

const (
	contextKeyQuerier = iota
	contextKeyUser
)

// WithQuerierInContext self explanatory!
func WithQuerierInContext(q *db.Queries, r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), key(contextKeyQuerier), q)
	return r.WithContext(ctx)
}

// GetQuerierFromContext self explanatory!
func GetQuerierFromContext(r *http.Request) *db.Queries {
	return r.Context().Value(key(contextKeyQuerier)).(*db.Queries)
}

func ise(w http.ResponseWriter, err error) {
	log.Printf("%v", err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "internal server error")
}

func forb(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, "forbidden")
}

func rj(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		ise(w, err)
	}
}

func rh(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, s)
}
