package admin

import (
	"log"
	"net/http"
	"yus/internal/handlers/common/response"
)

func (h *AdminHandler) GetBackupRoutes(w http.ResponseWriter, r *http.Request) {
	backup_routes, err := h.Store.DB.GetBackupRoutes(r.Context())
	if err != nil {
		log.Println("error while accessing backuproutes - ", err)
	}
	log.Println("backuproutes - ", backup_routes)
	response.WriteJSON(w, r, backup_routes)
}
