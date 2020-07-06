package middlewares

import (
	"log"
	"net/http"
	"strings"
	gdb "td4/back/db/generated"
	"td4/back/server_api/handlers"
	"time"

	"github.com/danilopolani/gocialite"
)

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

// Logging a logging middleware
func Logging(next http.Handler, q *gdb.Queries, g *gocialite.Dispatcher, corsOrigin string) http.Handler {
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
