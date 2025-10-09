package handlers

import (
	"net/http"
	"yus/internal/storage/postgres"
)

func Get_Schedule_handler(w http.ResponseWriter, r *http.Request) {
	// if !FindAdminSession(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	current_schedule:=postgres.Get_Current_schedule()
	WriteJSON(w,r,current_schedule)
}
