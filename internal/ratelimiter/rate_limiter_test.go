package ratelimiter

import (
	"context"
	"testing"

	"github.com/evg555/antibrutforce/internal/config"
)

func TestAAllowAttempt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	login := "login"
	password := "password"
	ip := "127.0.0.1"

	limit := 10
	expiration := 2

	rateLimiter := NewAuthRateLimiter(ctx, config.RateLimiter{
		LoginLimit:         limit,
		PasswordLimit:      limit * 10,
		IPLimit:            limit * 100,
		ExpirationInterval: expiration,
	})

	t.Run("first 10 requests", func(t *testing.T) {
		for i := 0; i < limit; i++ {
			if !rateLimiter.AllowAttempt(login, password, ip) {
				t.Errorf("request %d must be allowed", i+1)
			}
		}
	})

	t.Run("next 11 request", func(t *testing.T) {
		if rateLimiter.AllowAttempt(login, password, ip) {
			t.Errorf("11 request must be not allowed")
		}
	})
}

func TestResetBucket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	login := "login"
	password := "password"
	ip := "127.0.0.1"

	limit := 10

	rateLimiter := NewAuthRateLimiter(ctx, config.RateLimiter{
		LoginLimit:    limit * 10,
		PasswordLimit: limit,
	})

	t.Run("first 10 requests", func(t *testing.T) {
		for i := 0; i < limit; i++ {
			if !rateLimiter.AllowAttempt(login, password, ip) {
				t.Errorf("request %d must be allowed", i+1)
			}
		}
	})

	t.Run("next 11 request", func(t *testing.T) {
		if rateLimiter.AllowAttempt(login, password, ip) {
			t.Errorf("11 request must be not allowed")
		}
	})

	rateLimiter.ResetBucket(password, ip)

	t.Run("next 10 requests after reset", func(t *testing.T) {
		for i := 0; i < limit; i++ {
			if !rateLimiter.AllowAttempt(login, password, ip) {
				t.Errorf("request %d must be allowed", i+1)
			}
		}
	})
}
