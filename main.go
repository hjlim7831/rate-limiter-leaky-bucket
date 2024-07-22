package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "https://httpbin.org"
)

func (lb *LeakyBucket) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lb.Add() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
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
	response := sendRequest(baseURL)
	if response != nil {
		body, _ := io.ReadAll(response.Body)
		w.WriteHeader(response.StatusCode)
		_, err := w.Write(body)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

}

func sendRequest(URL string) *http.Response {
	log.Printf("Sending request to %s\n", URL)
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal("Error sending request: ", err)
		return nil
	}
	return resp
}
