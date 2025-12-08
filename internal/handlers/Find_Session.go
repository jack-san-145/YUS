package handlers

import (
	"net/http"
	"yus/internal/storage/redis"
)

func FindAdminSession_mobile(r *http.Request) bool {
	session_id := r.Header.Get("Authorization")
	if session_id == "" {
		return false
	}

	if redis.Check_Admin_session(session_id) {
		return true
	}
	return false
}

func FindAdminSession_web(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}

	if redis.Check_Admin_session(cookie.Value) {
		return true
	}
	return false
}

func FindDriver_httpSession(r *http.Request) (bool, int) {

	ctx := r.Context()
	session_id := r.Header.Get("Authorization")
	if session_id == "" {
		return false, 0
	}

	is_valid, driver_id, _ := redis.CheckDriverSession(ctx, session_id)
	if is_valid {
		return true, driver_id
	}
	return false, 0
}

func FindDriver_wssSession(r *http.Request) (bool, int) {
	ctx := r.Context()
	session_id := r.URL.Query().Get("session_id")
	if session_id == "" {
		return false, 0
	}

	is_valid, driver_id, _ := redis.CheckDriverSession(ctx, session_id)
	if is_valid {
		return true, driver_id
	}
	return false, 0
}
