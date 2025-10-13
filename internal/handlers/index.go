package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Serve_admin_login(w http.ResponseWriter, r *http.Request) {

	// if !FindAdminSession(r) {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	templ, err := template.ParseFiles("../../ui/templates/logo.html")
	if err != nil {
		fmt.Println("logo.html not found - ", err)
		return
	}

	// err=templ.Execute(w, nil)
	templ.Execute(w, nil)

}
