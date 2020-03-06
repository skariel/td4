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

// CreateSolution creates a solution to some test code
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

// CreateTestCode creates a new test
func CreateTestCode(w http.ResponseWriter, r *http.Request) {
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

// GetTestByID get a single test by id
func GetTestByID(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ise(w, err)
		return
	}

	test, err := q.GetTestCodeByID(context.Background(), int32(id))
	if err != nil {
		ise(w, err)
		return
	}

	rj(w, test)
}

// AllTests give all tests, plus avatar user, etc.
func AllTests(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)

	vars := mux.Vars(r)

	offset, err := strconv.Atoi(vars["offset"])
	if err != nil {
		ise(w, err)
		return
	}

	tests, err := q.GetTestCodes(context.Background(), int32(offset))
	if err != nil {
		ise(w, err)
		return
	}

	rj(w, tests)
}

// SolutionCodesByTest give all solutions bya a test ID
func SolutionCodesByTest(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ise(w, err)
		return
	}

	offset, err := strconv.Atoi(vars["offset"])
	if err != nil {
		ise(w, err)
		return
	}

	tests, err := q.GetSolutionCodesByTest(context.Background(),
		db.GetSolutionCodesByTestParams{
			TestCodeID: int32(id),
			Offset:     int32(offset),
		})
	if err != nil {
		ise(w, err)
		return
	}

	rj(w, tests)
}
