package main

import (
    "log"
    "net/http"
    pok "pok/src"
)

func main() {
    db, err := pok.InitDB("./pok.db")
    if err != nil {
        log.Fatalf("Échec de l'initialisation de la base de données : %v", err)
    }
    defer db.Close()

    // Force a fresh login after each server restart.
    if err := pok.DeleteAllSessions(); err != nil {
        log.Printf("Impossible de vider les sessions au démarrage : %v", err)
    }

    // Auth
    http.HandleFunc("/", pok.HomeHandler)
    http.HandleFunc("/login", pok.RedirectIfAuthenticated(pok.LoginHandler))
    http.HandleFunc("/connexion", pok.RedirectIfAuthenticated(pok.LoginHandler))
    http.HandleFunc("/register", pok.RedirectIfAuthenticated(pok.RegisterHandler))
    http.HandleFunc("/inscription", pok.RedirectIfAuthenticated(pok.RegisterHandler))

    // Protected (pour l’instant: juste dashboard + logout)
    http.HandleFunc("/dashboard", pok.DashboardHandler)
    http.HandleFunc("/logout", pok.LogoutHandler)

    // Static
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    log.Fatal(http.ListenAndServe(":8080", nil))
}