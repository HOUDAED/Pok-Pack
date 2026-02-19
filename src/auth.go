package pok

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginViewData struct {
	Error   string
	Success string
}

type RegisterViewData struct {
	Error string
}

func renderStatic(w http.ResponseWriter, filename string, data any) {
	t := template.Must(template.ParseFiles("static/" + filename))
	_ = t.Execute(w, data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
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
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := LoginViewData{}
		if r.URL.Query().Get("created") == "1" {
			data.Success = "Compte créé avec succès. Vous pouvez vous connecter."
		}
		if r.URL.Query().Get("error") == "1" {
			data.Error = "Identifiants incorrects."
		}
		renderStatic(w, "connexion.html", data)
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}
		login := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		if login == "" || password == "" {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}

		u, err := GetUserByPseudoOrEmail(login)
		if err != nil {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}
		if !CheckPasswordHash(password, u.PasswordHash) {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}

		sid, err := newSessionID()
		if err != nil {
			log.Printf("Erreur génération session: %v", err)
			http.Error(w, "Erreur serveur.", http.StatusInternalServerError)
			return
		}
		now := time.Now()
		expires := now.Add(14 * 24 * time.Hour)
		if err := InsertSession(sid, u.ID, expires.Unix(), now.Unix()); err != nil {
			log.Printf("Erreur insertion session: %v", err)
			http.Error(w, "Erreur serveur.", http.StatusInternalServerError)
			return
		}
		setSessionCookie(w, sid, expires)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}


func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderStatic(w, "inscription.html", RegisterViewData{})
		return
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			renderStatic(w, "inscription.html", RegisterViewData{Error: "Formulaire invalide."})
			return
		}

		pseudo := strings.TrimSpace(r.FormValue("username"))
		if pseudo == "" {
			pseudo = strings.TrimSpace(r.FormValue("pseudo"))
		}
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		if pseudo == "" || email == "" || password == "" {
			renderStatic(w, "inscription.html", RegisterViewData{Error: "Merci de remplir tous les champs."})
			return
		}
		if IsEmailTaken(email) {
			renderStatic(w, "inscription.html", RegisterViewData{Error: "Cet email est déjà utilisé."})
			return
		}
		if IsPseudoTaken(pseudo) {
			renderStatic(w, "inscription.html", RegisterViewData{Error: "Ce pseudo est déjà pris."})
			return
		}
		if !IsPasswordValid(password) {
			renderStatic(w, "inscription.html", RegisterViewData{Error: "Mot de passe non conforme aux règles CNIL."})
			return
		}
		if err := CreateUser(pseudo, email, password); err != nil {
			log.Printf("Erreur création utilisateur : %v", err)
			renderStatic(w, "inscription.html", RegisterViewData{Error: "Erreur lors de la création du compte."})
			return
		}

		http.Redirect(w, r, "/login?created=1", http.StatusSeeOther)
		return
	default:
		http.Redirect(w, r, "/inscription", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sid := getSessionIDFromRequest(r)
	if sid != "" {
		_ = DeleteSession(sid)
	}
	clearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
