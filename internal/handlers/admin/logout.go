package admin

import (
	"log"
	"net/http"
	"yus/internal/handlers/common/response"
)

func (h *AdminHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session_id := r.Context().Value("ADMIN_SESSION").(string)
	log.Println("admin session - ", session_id)
	err := h.Store.InMemoryDB.DeleteSession(ctx, session_id)
	if err != nil {
		response.WriteJSON(w, r, map[string]bool{"status": false})
		return
	}
	// response.WriteJSON(w, r, map[string]bool{"status": true})
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	response.WriteJSON(w, r, map[string]bool{
		"status": true,
	})
}
