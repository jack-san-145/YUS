package admin

import (
	"net/http"
	"yus/internal/handlers/common/response"
)

func (h *AdminHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session_id := r.Context().Value("ADMIN_SESSION").(string)

	err := h.Store.InMemoryDB.RemoveAdminSession(ctx, session_id)
	if err != nil {
		response.WriteJSON(w, r, map[string]bool{"status": false})
		return
	}
	// http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
