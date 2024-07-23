package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLeakyBucket_Add(t *testing.T) {
	leakyBucket := NewLeakyBucket(2, 10*time.Millisecond)

	req1 := Request{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "http://example.com", nil), done: make(chan struct{})}
	req2 := Request{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "http://example.com", nil), done: make(chan struct{})}
	req3 := Request{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "http://example.com", nil), done: make(chan struct{})}

	if !leakyBucket.Add(req1) {
		t.Fatal("expected request 1 to be added to the bucket")
	}

	if !leakyBucket.Add(req2) {
		t.Fatal("expected request 2 to be added to the bucket")
	}

	if leakyBucket.Add(req3) {
		t.Fatal("expected request 3 to be rejected from the bucket")
	}
}

func TestLeakyBucket_HandleRequest(t *testing.T) {
	leakyBucket := NewLeakyBucket(1, 10*time.Millisecond)

	req := Request{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "http://example.com", nil), done: make(chan struct{})}

	if !leakyBucket.Add(req) {
		t.Fatal("expected request to be added to the bucket")
	}

	<-req.done // 대기하여 요청이 완료되었는지 확인

	// 요청이 완료된 후 응답 확인
	recorder := req.w.(*httptest.ResponseRecorder)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", recorder.Code)
	}
}

func TestLeakyBucket_StartLeaking(t *testing.T) {
	leakyBucket := NewLeakyBucket(1, 10*time.Millisecond)

	req := Request{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "http://example.com", nil), done: make(chan struct{})}

	if !leakyBucket.Add(req) {
		t.Fatal("expected request to be added to the bucket")
	}

	<-req.done // 대기하여 요청이 완료되었는지 확인

	// 요청이 완료된 후 응답 확인
	recorder := req.w.(*httptest.ResponseRecorder)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", recorder.Code)
	}
}

func TestLeakyBucket_EmptyBucket(t *testing.T) {
	leakyBucket := NewLeakyBucket(1, 10*time.Millisecond)

	time.Sleep(20 * time.Millisecond) // 버킷이 비어 있는 상태를 만들기 위해 대기

	req := Request{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "http://example.com", nil), done: make(chan struct{})}

	if !leakyBucket.Add(req) {
		t.Fatal("expected request to be added to the bucket")
	}

	<-req.done // 대기하여 요청이 완료되었는지 확인

	// 요청이 완료된 후 응답 확인
	recorder := req.w.(*httptest.ResponseRecorder)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", recorder.Code)
	}
}
