package handlers

import (
	"net/http"
	"yus/internal/storage/redis"
)

func FindSession(r *http.Request) bool {
	session_id := r.Header.Get("Authorization")
	if session_id == "" {
		return false
	}

	if redis.Check_session(session_id) {
		return true
	}
	return false
}
