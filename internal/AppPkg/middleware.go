package AppPkg

import (
	"context"
	"fmt"
	"net/http"
)

func (app *Application) IsDriverAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var sessionID string
		sessionID = r.Header.Get("Authorization")
		if sessionID == "" {
			sessionID = r.URL.Query().Get("session_id")
		}
		if sessionID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		valid, driverID, _ := app.Store.InMemoryDB.CheckDriverSession(r.Context(), sessionID)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "DRIVER_ID", driverID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (app *Application) IsAdminAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Try getting sessionID from header first
		sessionID := r.Header.Get("Authorization")

		// If header is empty, try cookie
		if sessionID == "" {
			cookie, err := r.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			sessionID = cookie.Value //assin cookie value as sessionID
		}

		// Check session validity
		valid, _ := app.Store.InMemoryDB.CheckAdminSession(ctx, sessionID)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println("Admin session id - ", sessionID)
		// Call next handler
		ctx = context.WithValue(ctx, "ADMIN_SESSION", sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
