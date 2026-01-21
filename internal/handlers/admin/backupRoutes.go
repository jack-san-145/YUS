package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"yus/internal/handlers/common/response"
	"yus/internal/models"
)

func (h *AdminHandler) GetBackupRoutesHandler(w http.ResponseWriter, r *http.Request) {
	backup_routes, err := h.Store.DB.GetBackupRoutes(r.Context())
	if err != nil {
		log.Println("error while accessing backuproutes - ", err)
	}
	log.Println("backuproutes - ", backup_routes)
	response.WriteJSON(w, r, backup_routes)
}

func (h *AdminHandler) SaveBackupRoutesHandler(w http.ResponseWriter, r *http.Request) {

	var NewRoute models.BackupRoute
	err := json.NewDecoder(r.Body).Decode(&NewRoute)
	if err != nil {
		log.Println("error while get new route - ", err)
		return
	}
	err = h.Store.DB.StoreFromBackupRoute(r.Context(), &NewRoute)
	if err != nil {
		response.WriteJSON(w, r, map[string]bool{"status": false})
		return
	}
	response.WriteJSON(w, r, map[string]bool{"status": true})
}
