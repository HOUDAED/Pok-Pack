package main

import (
	"log"
	"net/http"
	pok "pok/src"
)

func main() {
	http.HandleFunc("/", pok.HomeHandler)
	http.HandleFunc("/login", pok.RedirectIfAuthenticated(pok.LoginHandler))
	http.HandleFunc("/connexion", pok.RedirectIfAuthenticated(pok.LoginHandler))
	http.HandleFunc("/register", pok.RedirectIfAuthenticated(pok.RegisterHandler))
	http.HandleFunc("/inscription", pok.RedirectIfAuthenticated(pok.RegisterHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
