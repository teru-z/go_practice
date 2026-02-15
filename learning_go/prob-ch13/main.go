package main

import (
	"fmt"
	"net/http"
	"time"
)

func getCurrentTimeInRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

func run() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprint(w, getCurrentTimeInRFC3339())
		}
	})
	http.ListenAndServe(":8080", nil)
}

func main() {
	run()
}
