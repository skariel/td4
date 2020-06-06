package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danilopolani/gocialite"

	"github.com/gorilla/mux"

	cache "github.com/victorspringer/http-cache"

	"github.com/victorspringer/http-cache/adapter/memory"

	"back/handlers"

	"sql/db"
)

// some conf
const (
	httptimeoutSeconds                  = 3.0
	cleaningPendingRunsPerUSerEveryrHrs = 8.0
	cleaningLongRunsEveryrHrs           = 1.0
	port                                = ":8081"
	corsOrigin                          = "*"
	maxTitleLen                         = 256
	maxDescLen                          = 2048
	maxCodeLen                          = 8192
	cacheCapacity                       = 50000
	cacheTTLSeconds                     = 7
)

func main() {
	// configure root directory
	td4Root := os.Getenv("TD4_ROOT")
	httpRoot := td4Root + "/back"
	log.Printf("httpRoot: %v", httpRoot)

	// connect to the DB
	q, _, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")

	// routing
	r := mux.NewRouter()

	// static files
	r.PathPrefix("/static/").Handler(
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir(httpRoot+"/static/")))).Methods("GET")

	// social login
	r.HandleFunc("/auth/github", handlers.SocialRedirectHandler).Methods("GET")
	r.HandleFunc("/auth/github/callback", handlers.SocialCallbackHandler).Methods("GET")

	// custom handlers
	r.HandleFunc("/api/create_test", handlers.CreateTestCodeConfigurator(maxTitleLen, maxDescLen, maxCodeLen)).Methods("POST")
	r.HandleFunc("/api/create_solution", handlers.CreateSolutionCodeConfigurator(maxCodeLen)).Methods("POST")
	r.HandleFunc("/api/test/{id}", handlers.GetTestByID).Methods("GET")
	r.HandleFunc("/api/alltests/{offset}", handlers.AllTests).Methods("GET")
	r.HandleFunc("/api/solutions_by_test/{id}/{offset}", handlers.SolutionCodesByTest).Methods("GET")
	r.HandleFunc("/api/solution/{id}", handlers.SolutionCodeByID).Methods("GET")
	r.HandleFunc("/api/results_by_run/{id}", handlers.ResultsByRun).Methods("GET")
	r.HandleFunc("/api/update_solution", handlers.UpdateSolutionCodeConfigurator(maxCodeLen)).Methods("POST")
	r.HandleFunc("/api/update_test", handlers.UpdateTestCodeConfigurator(maxTitleLen, maxDescLen, maxCodeLen)).Methods("POST")

	// apply middlewares

	// timeout
	h := http.TimeoutHandler(r, httptimeoutSeconds*time.Second, "Timeout!\n")

	// my middleware (cors, logging, context etc.)
	h = middleware(h, q, gocialite.NewDispatcher())

	// caching. doesn't reach logging
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(cacheCapacity),
	)
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(cacheTTLSeconds*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	h = cacheClient.Middleware(h)

	// start some cleanup functions
	go doEvery(cleaningPendingRunsPerUSerEveryrHrs*time.Hour, func() {
		log.Println("cleaning pending runs per user")
		err = q.CleanPendingRunsPerUSer(context.Background())
		if err != nil {
			log.Printf("error while cleaning pending tasks per use: %v", err)
		}
	})
	go doEvery(cleaningLongRunsEveryrHrs*time.Hour, func() {
		log.Println("cleaning long runs")
		err = q.FailLongRuns(context.Background())
		if err != nil {
			log.Printf("error while cleaning long runs: %v", err)
		}
	})

	// serve!
	log.Println("Serving at " + port)
	log.Fatal(http.ListenAndServeTLS(port, httpRoot+"/server.crt", httpRoot+"/server.key", h))
}

func doEvery(d time.Duration, fn func()) {
	fn()

	for range time.Tick(d) {
		fn()
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}

	n, err := w.ResponseWriter.Write(b)
	w.length += n

	return n, err
}

func middleware(next http.Handler, q *db.Queries, g *gocialite.Dispatcher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// enable CORS
		w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// preflight stuff...
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// get user and put into context
		user := handlers.GetUserFromAuthorizationHeader(r)
		r = handlers.WithUserInContext(user, r)

		// put querier into context
		r = handlers.WithQuerierInContext(q, r)

		// put gocial into context
		r = handlers.WithGocialInContext(g, r)

		// Trim slash
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		// Next!
		sw := statusWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r)

		// Log
		displayName := "nil"
		if user != nil {
			displayName = user.DisplayName
		}
		log.Printf("%v %v %v %v %v %v %v %v",
			time.Since(startTime),
			r.RemoteAddr,
			r.Proto, r.Method,
			r.RequestURI,
			sw.status,
			sw.length,
			displayName)
	})
}
