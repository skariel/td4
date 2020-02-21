package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"../sql/db"
	"./handlers"
)

var (
	q *db.Queries
)

func main() {
	td4Root := os.Getenv("TD4_ROOT")
	if td4Root == "" {
		td4Root = "."
		log.Printf("TD4_ROOT env is missing. Using working path as default")
	}
	httpRoot := td4Root + "/back"
	// connect to the DB
	q, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB")

	// routing
	handlers.SetQueries(q)
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
	h := http.TimeoutHandler(r, 3*time.Second, "Timeout!\n")
	h = middleware(h)

	// start some cleanup functions
	go doEvery(8*time.Hour, func() {
		log.Println("cleaning pending runs per user")
		q.CleanPendingRunsPerUSer(context.Background())
	})

	// serve!
	log.Println("Serving at localhost:8081")
	log.Fatal(http.ListenAndServeTLS(":8081", httpRoot+"/server.crt", httpRoot+"/server.key", h))
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

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		// enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// preflight stuff...
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// get user
		user := handlers.GetUserFromAuthorizationHeader(r)
		r = handlers.WithUserInContext(user, r)

		// Trim slash
		if len(r.URL.Path) > 1 {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		// Next!
		sw := statusWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r)

		// Log
		log.Printf("%v %v %v %v %v %v %v", time.Since(startTime), r.RemoteAddr, r.Proto, r.Method, r.RequestURI, sw.status, sw.length)
	})
}
