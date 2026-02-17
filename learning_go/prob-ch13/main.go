package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func getCurrentTimeInRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("request received",
				"method", r.Method,
				"path", r.URL.Path,
				"remote", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
			next.ServeHTTP(w, r)
		})
	}
}

func run() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprint(w, getCurrentTimeInRFC3339())
		}
	})

	handler := LoggingMiddleware(logger)(mux)
	http.ListenAndServe(":8080", handler)
}

func main() {
	run()
}
