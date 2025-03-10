package rate_limiter

import (
	"sync"
	"time"
)

type Bucket struct {
	capacity  int
	rate      time.Duration
	lastCheck time.Time
	mutex     sync.Mutex
	tokens    int
}

func NewBucket(capacity int, rate time.Duration) *Bucket {
	return &Bucket{
		capacity:  capacity,
		rate:      rate,
		lastCheck: time.Now(),
		tokens:    capacity,
	}
}

func (b *Bucket) Allow() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	now := time.Now()
	delta := now.Sub(b.lastCheck)
	tokensToAdd := int(delta / b.rate)

	if tokensToAdd > 0 {
		b.tokens = min(b.capacity, b.tokens+tokensToAdd)
		b.lastCheck = now
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}
