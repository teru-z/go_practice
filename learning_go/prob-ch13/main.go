package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func getCurrentTimeInRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

type NowJSON struct {
	DayOfWeek  string `json:"day_of_week"`
	DayOfMonth int    `json:"day_of_month"`
	Month      int    `json:"month"`
	Year       int    `json:"year"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
	Second     int    `json:"second"`
}

func getCurrentTimeInJSONString() string {
	now := time.Now()

	v := NowJSON{
		DayOfWeek:  now.Weekday().String(), // "Monday" など
		DayOfMonth: now.Day(),
		Month:      int(now.Month()),
		Year:       now.Year(),
		Hour:       now.Hour(),
		Minute:     now.Minute(),
		Second:     now.Second(),
	}

	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("request received",
				"method", r.Method,
				"header", r.Header,
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
			if r.Header.Get("Accept") == "application/json" {
				fmt.Fprint(w, getCurrentTimeInJSONString())
			} else {
				fmt.Fprint(w, getCurrentTimeInRFC3339())
			}
		}
	})

	handler := LoggingMiddleware(logger)(mux)
	http.ListenAndServe(":8080", handler)
}

func main() {
	run()
}
