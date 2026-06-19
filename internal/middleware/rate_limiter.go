package middleware

import (
"sync"
"time"

"github.com/gin-gonic/gin"
)

type RateLimiter struct {
visitors map[string]*Visitor
mu       sync.Mutex
limit    int
duration time.Duration
}

type Visitor struct {
count     int
lastSeen  time.Time
resetTime time.Time
}

func NewRateLimiter(limit int, duration time.Duration) *RateLimiter {
rl := &RateLimiter{
visitors: make(map[string]*Visitor),
limit:    limit,
duration: duration,
}

// Cleanup expired entries
go rl.cleanup()

return rl
}

func (rl *RateLimiter) cleanup() {
for {
time.Sleep(time.Minute)
rl.mu.Lock()
for ip, visitor := range rl.visitors {
if time.Now().After(visitor.resetTime) {
delete(rl.visitors, ip)
}
}
rl.mu.Unlock()
}
}

func (rl *RateLimiter) Allow(ip string) bool {
rl.mu.Lock()
defer rl.mu.Unlock()

visitor, exists := rl.visitors[ip]
if !exists {
rl.visitors[ip] = &Visitor{
count:     1,
lastSeen:  time.Now(),
resetTime: time.Now().Add(rl.duration),
}
return true
}

if time.Now().After(visitor.resetTime) {
visitor.count = 1
visitor.resetTime = time.Now().Add(rl.duration)
visitor.lastSeen = time.Now()
return true
}

if visitor.count >= rl.limit {
return false
}

visitor.count++
visitor.lastSeen = time.Now()
return true
}

func RateLimitMiddleware(limit int, duration time.Duration) gin.HandlerFunc {
limiter := NewRateLimiter(limit, duration)

return func(c *gin.Context) {
ip := c.ClientIP()
if !limiter.Allow(ip) {
c.JSON(429, gin.H{
"error": "Rate limit exceeded. Please try again later.",
"limit": limit,
"duration": duration.String(),
})
c.Abort()
return
}
c.Next()
}
}
