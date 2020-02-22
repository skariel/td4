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

	"../sql/db"
	"./handlers"
)

// some conf
const (
	httptimeoutSeconds                  = 3.0
	cleaningPendingRunsPerUSerEveryrHrs = 8.0
	port                                = ":8081"
	corsOrigin                          = "*"
)

func main() {
	// configure root directory
	td4Root := os.Getenv("TD4_ROOT")
	httpRoot := td4Root + "/back"
	log.Printf("httpRoot: %v", httpRoot)

	// connect to the DB
	q, err := db.ConnectDB()
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
	r.HandleFunc("/api/tests", handlers.AllTests).Methods("GET")
	r.HandleFunc("/api/test/{id}", handlers.TestByID).Methods("GET")
	r.HandleFunc("/api/test/{id}/codes", handlers.CodesByTest).Methods("GET")
	r.HandleFunc("/api/code/{id}", handlers.TestCodeByID).Methods("GET")
	r.HandleFunc("/api/code/{id}/solutions", handlers.SolutionsByCode).Methods("GET")
	r.HandleFunc("/api/users", handlers.Users).Methods("GET")
	r.HandleFunc("/api/create_test", handlers.CreateTest).Methods("POST")
	r.HandleFunc("/api/create_code", handlers.CreateCode).Methods("POST")
	r.HandleFunc("/api/create_solution", handlers.CreateSolution).Methods("POST")

	// apply middlewares
	h := http.TimeoutHandler(r, httptimeoutSeconds*time.Second, "Timeout!\n")
	h = middleware(h, q, gocialite.NewDispatcher())

	// start some cleanup functions
	go doEvery(cleaningPendingRunsPerUSerEveryrHrs*time.Hour, func() {
		log.Println("cleaning pending runs per user")
		err = q.CleanPendingRunsPerUSer(context.Background())
		if err != nil {
			log.Printf("error while cleaning pending tasks per use: %v", err)
		}
	})

	// serve!
	log.Println("Serving at " + port)
	log.Fatal(http.ListenAndServeTLS(port, httpRoot+"/server.crt", httpRoot+"/server.key", h))
}

func doEvery(d time.Duration, fn func()) {
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
		log.Printf("%v %v %v %v %v %v %v", time.Since(startTime), r.RemoteAddr, r.Proto, r.Method, r.RequestURI, sw.status, sw.length)
	})
}
