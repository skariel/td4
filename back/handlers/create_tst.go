package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"../../sql/db"
)

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
