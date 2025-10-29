package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"forum/helpers"
)

type RateLimitLogin struct {
	Count        int
	BlockedUntil time.Time
	FirstTime    time.Time
	UserIP       string
}

type RateLimiterManager struct {
	limits map[string]*RateLimitLogin
	mu     sync.Mutex
	Limit  int
	Window time.Duration
}

func NewRateLimiterManager(limit int, window time.Duration) *RateLimiterManager {
	return &RateLimiterManager{
		limits: make(map[string]*RateLimitLogin),
		Limit:  limit,
		Window: window,
	}
}

func (m *RateLimiterManager) Check(ip string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	rl, exists := m.limits[ip]

	if !exists {

		m.limits[ip] = &RateLimitLogin{
			Count:        1,
			FirstTime:    now,
			BlockedUntil: time.Time{},
			UserIP:       ip,
		}
		return true
	}

	if now.Before(rl.BlockedUntil) {
		return false
	}

	if now.After(rl.BlockedUntil) && rl.Count > m.Limit {
		rl.Count = 0
		rl.FirstTime = now
		rl.BlockedUntil = time.Time{}
	}

	rl.Count++

	if rl.Count > m.Limit {
		rl.BlockedUntil = now.Add(m.Window)
		return false
	}

	return true
}

func GetUserIP(r *http.Request) string {
	ipPort := r.RemoteAddr
	ip := strings.Split(ipPort, ":")[0]
	return ip
}

func RateLimitLoginMiddleware(manager *RateLimiterManager, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			ip := GetUserIP(r)

			if !manager.Check(ip) {
				helpers.Errorhandler(w, "To many requests slow down", http.StatusTooManyRequests)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}
