package pok

import "net/http"

func ConnexionHandler(w http.ResponseWriter, r *http.Request) {
	RedirectIfAuthenticated(LoginHandler)(w, r)
}

func InscriptionHandler(w http.ResponseWriter, r *http.Request) {
	RedirectIfAuthenticated(RegisterHandler)(w, r)
}

