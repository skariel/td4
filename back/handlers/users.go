// Package handlers all http api endpoints implementations
package handlers

import (
	"context"
	"net/http"
)

// Users display a user list of tests
func Users(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)
	tests, err := q.GetUsersByID(context.Background())

	if err != nil {
		ise(w, err)
		return
	}

	rj(w, tests)
}
