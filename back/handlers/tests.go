package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// AllTests display a user list of tests
func AllTests(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)

	tests, err := q.GetTestsByDate(context.Background())
	if err != nil {
		ise(w, err)
		return
	}
	rj(w, tests)
}

// TestByID get a test by id..
func TestByID(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ise(w, err)
		return
	}
	test, err := q.GetTestByID(context.Background(), int32(id))
	if err != nil {
		ise(w, err)
		return
	}
	rj(w, test)
}

// CodesByTest get all codes of test by test ID
func CodesByTest(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ise(w, err)
		return
	}
	codes, err := q.GetTestCodesByTest(context.Background(), int32(id))
	if err != nil {
		ise(w, err)
		return
	}
	rj(w, codes)
}

// TestCodeByID get test code by ID
func TestCodeByID(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ise(w, err)
		return
	}
	code, err := q.GetTestCodeByID(context.Background(), int32(id))
	if err != nil {
		ise(w, err)
		return
	}
	rj(w, code)
}

// SolutionsByCode get test code by ID
func SolutionsByCode(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ise(w, err)
		return
	}
	code, err := q.GetSolutionsByCode(context.Background(), int32(id))
	if err != nil {
		ise(w, err)
		return
	}
	rj(w, code)
}
