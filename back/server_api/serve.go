package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/danilopolani/gocialite"

	"github.com/gorilla/mux"

	"td4/back/db"
	"td4/back/server_api/handlers"
	"td4/back/server_api/middlewares"
	"td4/back/utils"
)

// some conf
const (
	httptimeout                     = 3.0 * time.Second
	cleaningPendingRunsPerUSerEvery = 8.0 * time.Hour
	cleaningLongRunsEvery           = 1.0 * time.Hour
	maxTitleLen                     = 256
	maxDescLen                      = 2048
	maxCodeLen                      = 8192
	cacheCapacity                   = 100000
	globalLimiterCleanEvery         = 120 * time.Second
	globalLimiterWindowSize         = 2 * time.Second
	globalLimiterMaxRate            = 4.0
	newTestLimiterCleanEvery        = 10 * time.Hour
	newTestLimiterWindowSize        = 5 * time.Hour
	newTestLimiterMaxRate           = 5.0
	editTestLimiterCleanEvery       = 2 * time.Hour
	editTestLimiterWindowSize       = 1 * time.Hour
	editTestLimiterMaxRate          = 10.0
)

func main() {
	port := ":" + utils.LoggedGetEnv("TD4_API_PORT")
	certificateFilePath := utils.LoggedGetEnv("TD4_CERTIFICATE_FILE_PATH")
	keyFilePath := utils.LoggedGetEnv("TD4_KEY_FILE_PATH")
	corsOrigin := utils.LoggedGetEnv("TD4_CORS_ORIGIN")
	clientID := utils.LoggedGetEnv("TD4_github_client_id")
	clientSecret := utils.LoggedGetEnv("TD4_github_client_secret")
	jwtSecret := []byte(utils.LoggedGetEnv("TD4_JWT_SECRET"))
	socialAuthFinalDest := utils.LoggedGetEnv("TD4_SOCIAL_AUTH_FINAL_DEST")
	socialRedirectURL := utils.LoggedGetEnv("TD4_SOCIAL_AUTH_REDIRECT")
	cacheTTLSeconds, err := strconv.ParseInt(utils.LoggedGetEnv("TD4_CACHE_TTL_SECONDS"), 10, 64)

	if err != nil {
		log.Fatal("bad TD4_CACHE_TTL_SECONDS env definition. Shoudld be able to transfor to an integer")
	}

	cacheTTL := time.Duration(cacheTTLSeconds * int64(time.Second))

	// connect to the DB
	q, _ := db.ConnectDB()

	log.Println("Connected to DB")

	// routing
	r := mux.NewRouter()

	// social login
	r.HandleFunc("/auth/github", handlers.CreateSocialRedirectHandlerConfigurator(clientID, clientSecret, socialRedirectURL)).Methods("GET")
	r.HandleFunc("/auth/github/callback", handlers.CreateSocialCallbackHandlerConfigurator(jwtSecret, socialAuthFinalDest)).Methods("GET")

	// custom handlers
	newTestLmt := middlewares.NewLimiter(newTestLimiterCleanEvery, newTestLimiterWindowSize, newTestLimiterMaxRate)
	editTestLmt := middlewares.NewLimiter(editTestLimiterCleanEvery, editTestLimiterWindowSize, editTestLimiterMaxRate)
	newSolutionLmt := middlewares.NewLimiter(newTestLimiterCleanEvery, newTestLimiterWindowSize, newTestLimiterMaxRate)
	editSolutionLmt := middlewares.NewLimiter(editTestLimiterCleanEvery, editTestLimiterWindowSize, editTestLimiterMaxRate)

	r.HandleFunc("/api/create_test",
		newTestLmt.Middleware(
			handlers.CreateTestCodeConfigurator(
				maxTitleLen, maxDescLen, maxCodeLen))).Methods("POST")

	r.HandleFunc("/api/create_solution",
		newSolutionLmt.Middleware(
			handlers.CreateSolutionCodeConfigurator(
				maxCodeLen))).Methods("POST")

	r.HandleFunc("/api/update_solution",
		editSolutionLmt.Middleware(
			handlers.UpdateSolutionCodeConfigurator(
				maxCodeLen))).Methods("POST")

	r.HandleFunc("/api/update_test",
		editTestLmt.Middleware(
			handlers.UpdateTestCodeConfigurator(
				maxTitleLen, maxDescLen, maxCodeLen))).Methods("POST")

	r.HandleFunc("/api/test/{id}", handlers.GetTestByID).Methods("GET")
	r.HandleFunc("/api/alltests/{offset}", handlers.AllTests).Methods("GET")
	r.HandleFunc("/api/alltests_by_user/{offset}/{displayname}", handlers.AllTestsByUser).Methods("GET")
	r.HandleFunc("/api/solutions_by_test/{id}/{offset}", handlers.SolutionCodesByTest).Methods("GET")
	r.HandleFunc("/api/solution/{id}", handlers.SolutionCodeByID).Methods("GET")
	r.HandleFunc("/api/results_by_run/{id}", handlers.ResultsByRun).Methods("GET")

	r.HandleFunc("/api/delete_test/{id}", handlers.DeleteTestByID).Methods("DELETE")
	r.HandleFunc("/api/delete_solution/{id}", handlers.DeleteSolutionByID).Methods("DELETE")

	// apply global middlewares
	h := http.TimeoutHandler(r, httptimeout, "Timeout!\n")
	lmt := middlewares.NewLimiter(globalLimiterCleanEvery, globalLimiterWindowSize, globalLimiterMaxRate)
	h = lmt.Handler(h)
	h = middlewares.Logging(h, q, gocialite.NewDispatcher(), corsOrigin, jwtSecret)
	cacheClient := middlewares.NewMemoryCache(cacheCapacity, cacheTTL)
	h = cacheClient.Middleware(h)

	// start some cleanup functions (DB)
	go utils.DoEvery(cleaningPendingRunsPerUSerEvery, func() {
		log.Println("cleaning pending runs per user")
		err := q.CleanPendingRunsPerUSer(context.Background())
		if err != nil {
			log.Printf("error while cleaning pending tasks per use: %v", err)
		}
	})
	go utils.DoEvery(cleaningLongRunsEvery, func() {
		log.Println("cleaning long runs")
		err := q.FailLongRuns(context.Background())
		if err != nil {
			log.Printf("error while cleaning long runs: %v", err)
		}
	})

	// serve!
	log.Printf("Process %v Serving at port %v\n", os.Getpid(), port)
	log.Fatal(http.ListenAndServeTLS(port, certificateFilePath, keyFilePath, h))
}
