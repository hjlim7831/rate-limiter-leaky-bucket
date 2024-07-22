package main

import (
	"time"
)

type LeakyBucket struct {
	capacity int           // 버킷 최대 용량
	rate     time.Duration // 누수율 (누수 간격)
	tokens   chan struct{} // 토큰을 저장하는 채널
}

func NewLeakyBucket(capacity int, rate time.Duration) *LeakyBucket {
	bucket := &LeakyBucket{
		capacity: capacity,
		rate:     rate,
		tokens:   make(chan struct{}, capacity),
	}

	go bucket.startLeaking()

	return bucket
}

func (lb *LeakyBucket) startLeaking() {
	ticker := time.NewTicker(lb.rate)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-lb.tokens:
			// 누수 성공
		default:
			// 버킷이 비어 있음
		}
	}
}

func (lb *LeakyBucket) Add() bool {
	select {
	case lb.tokens <- struct{}{}:
		// 토큰 추가 성공
		return true
	default:
		return false
	}
}
