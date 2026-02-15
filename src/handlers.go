package pok

import "net/http"

func ConnexionHandler(w http.ResponseWriter, r *http.Request) {
	RedirectIfAuthenticated(LoginHandler)(w, r)
}

func InscriptionHandler(w http.ResponseWriter, r *http.Request) {
	RedirectIfAuthenticated(RegisterHandler)(w, r)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	RequireAuth(func(w http.ResponseWriter, r *http.Request, u User) {
		http.ServeFile(w, r, "static/dashboard.html")
	},
	)(w, r)
}
