package ratelimiter

import (
	"testing"
	"time"
)

func TestAllow(t *testing.T) {
	limit := 10
	rate := 100 * time.Millisecond

	bucket := NewBucket(limit, rate)

	t.Run("first 10 requests", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if !bucket.Allow() {
				t.Errorf("request %d must be allowed", i+1)
			}
		}
	})

	t.Run("next 11 request", func(t *testing.T) {
		if bucket.Allow() {
			t.Errorf("11 request must be not allowed")
		}
	})

	time.Sleep(200 * time.Millisecond)

	t.Run("2 tokens restored", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			if !bucket.Allow() {
				t.Errorf("request %d must be allowed", i+1)
			}
		}
	})

	time.Sleep(1 * time.Second)

	t.Run("all tokens restored", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if !bucket.Allow() {
				t.Errorf("request %d must be allowed", i+1)
			}
		}
	})
}

func TestReset(t *testing.T) {
	limit := 10
	rate := 100 * time.Millisecond

	bucket := NewBucket(limit, rate)

	t.Run("first 10 requests", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if !bucket.Allow() {
				t.Errorf("request %d must be allowed", i+1)
			}
		}
	})

	bucket.Reset()

	t.Run("next 10 requests", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if !bucket.Allow() {
				t.Errorf("request %d must be allowed after reset", i+1)
			}
		}
	})
}
