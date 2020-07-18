package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "td4/back/db/generated"
)

type key int

const (
	contextKeyQuerier = iota
	contextKeyGocial
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

// Ise Internal Server Error
func Ise(w http.ResponseWriter, err error) {
	log.Printf("%v", err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "internal server error")
}

// Forb Forbidden
func Forb(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, "forbidden")
}

// Limited well, limited!
func Limited(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprint(w, "you have hit the request rate limit. Please try again later")
}

// ExpectationFailure expected something else
func ExpectationFailure(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusExpectationFailed)
	fmt.Fprint(w, err)
}

// Rj respond with json
func Rj(w http.ResponseWriter, i interface{}) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		Ise(w, err)
	}
}
