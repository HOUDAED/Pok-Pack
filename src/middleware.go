package pok

import (
	"log"
	"net/http"
)

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
