package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Request struct {
	w    http.ResponseWriter
	r    *http.Request
	done chan struct{} // 완료 신호를 위한 채널
}

type LeakyBucket struct {
	capacity int           // 버킷 최대 용량
	rate     time.Duration // 누수율 (누수 간격)
	queue    chan Request  // 토큰을 저장하는 채널
	mu       sync.Mutex
}

func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	bucket := &LeakyBucket{
		capacity: capacity,
		rate:     rate,
		queue:    make(chan Request, capacity),
	}

	go bucket.startLeaking()

	return bucket
}

func (lb *LeakyBucket) startLeaking() {
	ticker := time.NewTicker(lb.rate)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case req := <-lb.queue:
			lb.handleRequest(req)
			close(req.done)
			// 누수 성공

		default:
			// 버킷이 비어 있음
		}
	}
}

func (lb *LeakyBucket) handleRequest(req Request) {
	log.Printf("handleRequest %s\n", req.r.URL)
	response := sendRequest(baseURL)
	if response != nil {
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)

		req.w.WriteHeader(response.StatusCode)
		_, err := req.w.Write(body)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		req.w.WriteHeader(http.StatusInternalServerError)
		_, err := req.w.Write([]byte("Internal Server Error"))
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

func (lb *LeakyBucket) Add(req Request) bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	select {
	case lb.queue <- req:
		// 토큰 추가 성공
		return true
	default:
		return false
	}
}
