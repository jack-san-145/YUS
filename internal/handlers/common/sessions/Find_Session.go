package handlers

import (
	"net/http"
	"yus/internal/storage/redis"
)

func FindAdminSession_mobile(r *http.Request) bool {
	ctx := r.Context()

	session_id := r.Header.Get("Authorization")
	if session_id == "" {
		return false
	}

	valid, _ := redis.NewRedisClient().CheckAdminSession(ctx, session_id)

	if valid {
		return true
	}
	return false
}

func FindAdminSession_web(r *http.Request) bool {
	ctx := r.Context()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	valid, _ := redis.NewRedisClient().CheckAdminSession(ctx, cookie.Value)
	if valid {
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

	is_valid, driver_id, _ := redis.NewRedisClient().CheckDriverSession(ctx, session_id)
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

	is_valid, driver_id, _ := redis.NewRedisClient().CheckDriverSession(ctx, session_id)
	if is_valid {
		return true, driver_id
	}
	return false, 0
}
