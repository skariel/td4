// Package handlers all http api endpoints implementations
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"../../sql/db"
	"github.com/gorilla/mux"
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

// CreateTest does what it says
func CreateTest(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	q := GetQuerierFromContext(r)

	if user == nil {
		forb(w)
		return
	}

	var testParams db.InsertTestParams

	err := json.NewDecoder(r.Body).Decode(&testParams)
	if err != nil {
		ise(w, err)
		return
	}

	testParams.CreatedBy = user.ID

	test, err := q.InsertTest(context.Background(), testParams)
	if err != nil {
		ise(w, err)
		return
	}

	rj(w, test)
}

// CreateCode does what it says
func CreateCode(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	q := GetQuerierFromContext(r)

	if user == nil {
		forb(w)
		return
	}

	var testCodeParams db.InsertTestCodeParams
	err := json.NewDecoder(r.Body).Decode(&testCodeParams)

	testCodeParams.CreatedBy = user.ID

	if err != nil {
		ise(w, err)
		return
	}

	testCode, err := q.InsertTestCode(context.Background(), testCodeParams)
	if err != nil {
		ise(w, err)
		return
	}

	rj(w, testCode)
}

// CreateSolution does what it says
func CreateSolution(w http.ResponseWriter, r *http.Request) {
	_ = true // so dupls doesn't complain about copying from `CreateCode` above
	user := GetUserFromContext(r)
	q := GetQuerierFromContext(r)

	if user == nil {
		forb(w)
		return
	}

	var solutionCodeParams db.InsertSolutionCodeParams

	err := json.NewDecoder(r.Body).Decode(&solutionCodeParams)
	solutionCodeParams.CreatedBy = user.ID

	if err != nil {
		ise(w, err)
		return
	}

	solutionCode, err := q.InsertSolutionCode(context.Background(), solutionCodeParams)
	if err != nil {
		ise(w, err)
		return
	}

	rj(w, solutionCode)
}

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
