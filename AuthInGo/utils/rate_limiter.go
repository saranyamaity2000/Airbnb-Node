package utils

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiterManager struct {
	ipLimiters   map[string]*rate.Limiter
	userLimiters map[int]*rate.Limiter
	rateLimit    rate.Limit
	burst        int
}

func NewRateLimitManager(rl rate.Limit, burst int) *RateLimiterManager {
	return &RateLimiterManager{
		rateLimit:    rl,
		burst:        burst,
		ipLimiters:   make(map[string]*rate.Limiter),
		userLimiters: make(map[int]*rate.Limiter),
	}
}

func (m *RateLimiterManager) GetOrCreateIPLimiter(ip string) *rate.Limiter {
	if limiter, exists := m.ipLimiters[ip]; exists {
		fmt.Println("FOUND ip limiter with ip ", ip)
		return limiter
	}
	fmt.Println("CREATING ip limiter with ip ", ip)
	limiter := rate.NewLimiter(m.rateLimit, m.burst)
	m.ipLimiters[ip] = limiter
	return limiter
}

func (m *RateLimiterManager) GetOrCreateUserLimiter(userId int) *rate.Limiter {
	if limiter, exists := m.userLimiters[userId]; exists {
		fmt.Println("FOUND user limiter with id ", userId)
		return limiter
	}
	fmt.Println("CREATING user limiter with userId ", userId)
	limiter := rate.NewLimiter(m.rateLimit, m.burst)
	m.userLimiters[userId] = limiter
	return limiter
}

func (rlm *RateLimiterManager) UserRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userLimiter *rate.Limiter
		userClaims, isUserClaims := (r.Context().Value("userClaims")).(*UserClaims)
		if isUserClaims {
			userLimiter = rlm.GetOrCreateUserLimiter(int(userClaims.UserId))
		}
		if userLimiter != nil && !userLimiter.Allow() {
			WriteJsonErrorResponse(w, http.StatusTooManyRequests, "Too Many Requests", fmt.Errorf("User Rate Limit Exceeded"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rlm *RateLimiterManager) IPRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ipLimiter *rate.Limiter
		ipHost, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ipLimiter = rlm.GetOrCreateIPLimiter(ipHost)
		} else {
			fmt.Println("No IP Host Found")
		}
		if ipLimiter != nil && !ipLimiter.Allow() {
			WriteJsonErrorResponse(w, http.StatusTooManyRequests, "Too Many Requests", fmt.Errorf("IP Rate Limit Exceeded"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

var RateLimitManagerInstance *RateLimiterManager

func init() {
	RateLimitManagerInstance = NewRateLimitManager(rate.Every(10*time.Second), 5)
}
