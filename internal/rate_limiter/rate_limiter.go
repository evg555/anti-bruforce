package rate_limiter

import (
	"context"
	"sync"
	"time"

	"github.com/evg555/antibrutforce/internal/config"
)

var (
	loginLimit         = 10
	passwordLimit      = 100
	ipLimit            = 1000
	expirationInterval = 300 * time.Second
)

type AuthRateLimiter struct {
	logins    map[string]*Bucket
	passwords map[string]*Bucket
	ips       map[string]*Bucket
	mutex     sync.Mutex
}

func NewAuthRateLimiter(ctx context.Context, cfg config.RateLimiter) *AuthRateLimiter {
	rl := &AuthRateLimiter{
		logins:    make(map[string]*Bucket),
		passwords: make(map[string]*Bucket),
		ips:       make(map[string]*Bucket),
	}

	if cfg.LoginLimit > 0 {
		loginLimit = cfg.LoginLimit
	}

	if cfg.PasswordLimit > 0 {
		passwordLimit = cfg.PasswordLimit
	}

	if cfg.IpLimit > 0 {
		ipLimit = cfg.IpLimit
	}

	if cfg.ExpirationInterval > 0 {
		expirationInterval = time.Duration(cfg.ExpirationInterval) * time.Second
	}

	go rl.cleanupBuckets(ctx)

	return rl
}

func (r *AuthRateLimiter) AllowAttempt(login, password, ip string) bool {
	loginBucket := r.getBucket(login, loginLimit, r.logins)
	passwordBucket := r.getBucket(password, passwordLimit, r.passwords)
	ipBucket := r.getBucket(ip, ipLimit, r.ips)

	return loginBucket.Allow() && passwordBucket.Allow() && ipBucket.Allow()
}

func (r *AuthRateLimiter) getBucket(key string, limit int, bucketMap map[string]*Bucket) *Bucket {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if bucket, exists := bucketMap[key]; exists {
		return bucket
	}

	bucket := NewBucket(limit, time.Minute/time.Duration(limit))
	bucketMap[key] = bucket

	return bucket
}

func (r *AuthRateLimiter) ResetBucket(password, ip string) {
	passwordBucket := r.getBucket(password, passwordLimit, r.passwords)
	ipBucket := r.getBucket(ip, ipLimit, r.ips)

	passwordBucket.Reset()
	ipBucket.Reset()
}

func (r *AuthRateLimiter) cleanupBuckets(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.mutex.Lock()

			now := time.Now()

			cleanup := func(bucketMap map[string]*Bucket) {
				for key, bucket := range bucketMap {
					bucket.mutex.Lock()
					lastUsed := bucket.lastCheck
					bucket.mutex.Unlock()

					if now.Sub(lastUsed) > expirationInterval {
						delete(bucketMap, key)
					}
				}
			}

			cleanup(r.logins)
			cleanup(r.passwords)
			cleanup(r.ips)

			r.mutex.Unlock()
		}
	}
}
