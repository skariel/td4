// Package handlers all http api endpoints implementations
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "td4/back/db/generated"

	"github.com/gorilla/mux"
)

// CreateSolutionCodeConfigurator returns a configured CreateSolutionCode handler
func CreateSolutionCodeConfigurator(maxCodeLen int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = true // so dupls doesn't complain about copying from `CreateCode` above
		user := GetUserFromContext(r)
		q := GetQuerierFromContext(r)

		if user == nil {
			Forb(w)
			return
		}

		var solutionCodeParams db.InsertSolutionCodeParams
		if err := json.NewDecoder(r.Body).Decode(&solutionCodeParams); err != nil {
			Ise(w, err)
			return
		}

		if len(solutionCodeParams.Code) > maxCodeLen {
			ExpectationFailure(w, fmt.Sprintf("code too long (%v > %v)", len(solutionCodeParams.Code), maxCodeLen))
			return
		}

		if solutionCodeParams.Code == "" {
			ExpectationFailure(w, "no code")
			return
		}

		solutionCodeParams.CreatedBy = user.ID

		solutionCode, err := q.InsertSolutionCode(context.Background(), solutionCodeParams)
		if err != nil {
			Ise(w, err)
			return
		}

		Rj(w, solutionCode)
	}
}

// CreateTestCodeConfigurator returns a configured CreatetestCode handler
func CreateTestCodeConfigurator(maxTitleLen, maxDescLen, maxCodeLen int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		q := GetQuerierFromContext(r)

		if user == nil {
			Forb(w)
			return
		}

		var testCodeParams db.InsertTestCodeParams
		if err := json.NewDecoder(r.Body).Decode(&testCodeParams); err != nil {
			Ise(w, err)
			return
		}

		if len(testCodeParams.Title) > maxTitleLen {
			ExpectationFailure(w, fmt.Sprintf("title too long (%v > %v)", len(testCodeParams.Title), maxTitleLen))
			return
		}

		if testCodeParams.Title == "" {
			ExpectationFailure(w, "no title")
			return
		}

		if len(testCodeParams.Descr) > maxDescLen {
			ExpectationFailure(w, fmt.Sprintf("description too long (%v > %v)", len(testCodeParams.Descr), maxDescLen))
			return
		}

		if testCodeParams.Descr == "" {
			ExpectationFailure(w, "no description")
			return
		}

		if len(testCodeParams.Code) > maxCodeLen {
			ExpectationFailure(w, fmt.Sprintf("code too long (%v > %v)", len(testCodeParams.Code), maxCodeLen))
			return
		}

		if testCodeParams.Code == "" {
			ExpectationFailure(w, "no code")
			return
		}

		testCodeParams.CreatedBy = user.ID

		testCode, err := q.InsertTestCode(context.Background(), testCodeParams)
		if err != nil {
			Ise(w, err)
			return
		}

		Rj(w, testCode)
	}
}

// GetTestByID get a single test by id
func GetTestByID(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		Ise(w, err)
		return
	}

	test, err := q.GetTestCodeByID(context.Background(), int32(id))
	if err != nil {
		Ise(w, err)
		return
	}

	Rj(w, test)
}

// AllTests give all tests, plus avatar user, etc.
func AllTests(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)

	vars := mux.Vars(r)

	offset, err := strconv.ParseInt(vars["offset"], 10, 32)
	if err != nil {
		Ise(w, err)
		return
	}

	tests, err := q.GetTestCodes(context.Background(), int32(offset))
	if err != nil {
		Ise(w, err)
		return
	}

	Rj(w, tests)
}

// SolutionCodesByTest give all solutions by test ID
func SolutionCodesByTest(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		Ise(w, err)
		return
	}

	offset, err := strconv.ParseInt(vars["offset"], 10, 32)
	if err != nil {
		Ise(w, err)
		return
	}

	solutions, err := q.GetSolutionCodesByTest(context.Background(),
		db.GetSolutionCodesByTestParams{
			TestCodeID: int32(id),
			Offset:     int32(offset),
		})
	if err != nil {
		Ise(w, err)
		return
	}

	Rj(w, solutions)
}

// SolutionCodeByID give specific solution
func SolutionCodeByID(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		Ise(w, err)
		return
	}

	solution, err := q.GetSolutionCodeByID(context.Background(), int32(id))
	if err != nil {
		Ise(w, err)
		return
	}

	Rj(w, solution)
}

// ResultsByRun give all results by result ID
func ResultsByRun(w http.ResponseWriter, r *http.Request) {
	q := GetQuerierFromContext(r)

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		Ise(w, err)
		return
	}

	results, err := q.GetResultsByRun(context.Background(), int32(id))
	if err != nil {
		Ise(w, err)
		return
	}

	Rj(w, results)
}

// UpdateSolutionCodeConfigurator returns a configured UpdateSolutionCode handler
func UpdateSolutionCodeConfigurator(maxCodeLen int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		q := GetQuerierFromContext(r)

		if user == nil {
			Forb(w)
			return
		}

		// decodce request body
		var solutionCodeParams db.UpdateSolutionCodeParams
		if err := json.NewDecoder(r.Body).Decode(&solutionCodeParams); err != nil {
			Ise(w, err)
			return
		}

		// validate code
		if len(solutionCodeParams.Code) > maxCodeLen {
			ExpectationFailure(w, fmt.Sprintf("code too long (%v > %v)", len(solutionCodeParams.Code), maxCodeLen))
			return
		}

		if solutionCodeParams.Code == "" {
			ExpectationFailure(w, "no code")
			return
		}

		// validate user: compare user to solution owner (don't allow anyone else to update)
		solution, err := q.GetSolutionCodeByID(context.Background(), solutionCodeParams.ID)
		if err != nil {
			Ise(w, err)
			return
		}

		if user.ID != solution.CreatedBy {
			Forb(w)
			return
		}

		// do the update
		err = q.UpdateSolutionCode(context.Background(), solutionCodeParams)
		if err != nil {
			Ise(w, err)
			return
		}

		Rj(w, "")
	}
}

// UpdateTestCodeConfigurator returns a configured UpdateTestCode handler
func UpdateTestCodeConfigurator(maxTitleLen, maxDescLen, maxCodeLen int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		q := GetQuerierFromContext(r)

		if user == nil {
			Forb(w)
			return
		}

		// decodce request body
		var testCodeParams db.UpdateTestCodeParams
		if err := json.NewDecoder(r.Body).Decode(&testCodeParams); err != nil {
			Ise(w, err)
			return
		}

		// validate
		if len(testCodeParams.Title) > maxTitleLen {
			ExpectationFailure(w, fmt.Sprintf("title too long (%v > %v)", len(testCodeParams.Title), maxTitleLen))
			return
		}

		if testCodeParams.Title == "" {
			ExpectationFailure(w, "no title")
			return
		}

		if len(testCodeParams.Descr) > maxDescLen {
			ExpectationFailure(w, fmt.Sprintf("description too long (%v > %v)", len(testCodeParams.Descr), maxDescLen))
			return
		}

		if testCodeParams.Descr == "" {
			ExpectationFailure(w, "no description")
			return
		}

		if len(testCodeParams.Code) > maxCodeLen {
			ExpectationFailure(w, fmt.Sprintf("code too long (%v > %v)", len(testCodeParams.Code), maxCodeLen))
			return
		}

		if testCodeParams.Code == "" {
			ExpectationFailure(w, "no code")
			return
		}

		// validate user: compare user to solution owner (don't allow anyone else to update)
		test, err := q.GetTestCodeByID(context.Background(), testCodeParams.ID)
		if err != nil {
			Ise(w, err)
			return
		}

		if user.ID != test.CreatedBy {
			Forb(w)
			return
		}

		// do the update
		err = q.UpdateTestCode(context.Background(), testCodeParams)
		if err != nil {
			Ise(w, err)
			return
		}

		Rj(w, "")
	}
}

// DeleteTestByID self explanatory
func DeleteTestByID(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	q := GetQuerierFromContext(r)

	// get test id
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		Ise(w, err)
		return
	}

	// validate user: make sure we even have one :)
	if user == nil {
		Forb(w)
		return
	}

	// validate user: compare user to solution owner (don't allow anyone else to update)
	test, err := q.GetTestCodeByID(context.Background(), int32(id))
	if err != nil {
		Ise(w, err)
		return
	}

	if user.ID != test.CreatedBy {
		Forb(w)
		return
	}

	// delete!
	err = q.DeleteTestByID(context.Background(), test.ID)
	if err != nil {
		Ise(w, err)
		return
	}

	Rj(w, "")
}

// DeleteSolutionByID self explanatory
func DeleteSolutionByID(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	q := GetQuerierFromContext(r)

	// get solution id
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		Ise(w, err)
		return
	}

	// validate user: make sure we even have one :)
	if user == nil {
		Forb(w)
		return
	}

	// validate user: compare user to solution owner (don't allow anyone else to update)
	test, err := q.GetSolutionCodeByID(context.Background(), int32(id))
	if err != nil {
		Ise(w, err)
		return
	}

	if user.ID != test.CreatedBy {
		Forb(w)
		return
	}

	// delete!
	err = q.DeleteSolutionByID(context.Background(), test.ID)
	if err != nil {
		Ise(w, err)
		return
	}

	Rj(w, "")
}
