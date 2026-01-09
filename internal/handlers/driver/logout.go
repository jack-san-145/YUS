package driver

import (
	"net/http"
	"yus/internal/handlers/common/response"
)

func (h *DriverHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionID := r.Header.Get("Authorization")
	if sessionID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := h.Store.InMemoryDB.DeleteSession(ctx, sessionID)
	if err != nil {
		response.WriteJSON(w, r, map[string]bool{"status": false})
		return
	}

	response.WriteJSON(w, r, map[string]bool{
		"status": true,
	})
}
