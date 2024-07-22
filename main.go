package main

import (
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "https://httpbin.org"
)

func (lb *LeakyBucket) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lb.Add(w, r) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
	})
}

func main() {
	leakyBucket := NewLeakyBucket(3, 10*time.Second)
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/request",
		handleRequest)

	rateLimitedMux := leakyBucket.Limit(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      rateLimitedMux,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}

}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	log.Printf("handleRequest %s\n", req.URL)
}
