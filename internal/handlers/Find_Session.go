package handlers

import (
	"net/http"
	"yus/internal/storage/redis"
)

func FindAdminSession(r *http.Request) bool {
	session_id := r.Header.Get("Authorization")
	if session_id == "" {
		return false
	}

	if redis.Check_Admin_session(session_id) {
		return true
	}
	return false
}

func FindDriverSession(r *http.Request) (bool, int) {
	session_id := r.Header.Get("Authorization")
	if session_id == "" {
		return false, 0
	}

	is_valid, driver_id := redis.Check_Driver_session(session_id)
	if is_valid {
		return true, driver_id
	}
	return false, 0
}
