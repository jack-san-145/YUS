package AppPkg

import (
	"net"
	"net/http"
	"strings"
	"time"
	"yus/internal/models"
)

func (app *Application) RateLimit(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var rateLimit models.RateLimit
		ctx := r.Context()

		ip := "rate:ip:" + GetClientIP(r)
		rateLimit.Key = ip
		rateLimit.Capacity = 10
		rateLimit.RefillPerSecond = 1
		rateLimit.TimeStamp = time.Now().UnixMilli()

		allowed, err := app.Store.InMemoryDB.RateLimiter(ctx, &rateLimit)

		if err != nil {
			http.Error(w, "Rate Limit Error", http.StatusInternalServerError)
			return
		}

		if allowed == 0 {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)

	})

}

func GetClientIP(r *http.Request) string {

	// 1. X-Forwarded-For (can contain multiple IPs)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// 2. X-Real-IP
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return strings.TrimSpace(ip)
	}

	// 3. RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return r.RemoteAddr
}
