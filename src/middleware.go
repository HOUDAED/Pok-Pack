package pok

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

// currentUser returns the authenticated user from the session cookie.
// ok=false means unauthenticated.
func currentUser(r *http.Request) (u User, ok bool, err error) {
	sid := getSessionIDFromRequest(r)
	if sid == "" {
		return User{}, false, nil
	}
	u, err = GetUserBySessionID(sid, time.Now().Unix())
	if err == sql.ErrNoRows {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, err
	}
	return u, true, nil
}

// RequireAuth protects a handler: unauthenticated users are redirected to /login.
func RequireAuth(next func(http.ResponseWriter, *http.Request, User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok, err := currentUser(r)
		if err != nil {
			log.Printf("Erreur auth/session: %v", err)
			http.Error(w, "Erreur serveur.", http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r, u)
	}
}

// RedirectIfAuthenticated prevents logged-in users from seeing login/register pages.
func RedirectIfAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok, err := currentUser(r)
		if err != nil {
			log.Printf("Erreur auth/session: %v", err)
			http.Error(w, "Erreur serveur.", http.StatusInternalServerError)
			return
		}
		if ok {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
