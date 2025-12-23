package driver

import (
	"net/http"
	"yus/internal/handlers/common/response"
)

func (h *DriverHandler) RemoveAccountHandler(w http.ResponseWriter, r *http.Request) {
	driver_id := r.Context().Value("DRIVER_ID").(int)
	err := h.Store.DB.StoreDriverRemovalRequest(r.Context(), driver_id)
	if err != nil {
		response.WriteJSON(w, r, map[string]bool{"status": false})
		return
	}
	response.WriteJSON(w, r, map[string]bool{"status": true})
}
