package pok

import "net/http"

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    RequireAuth(func(w http.ResponseWriter, r *http.Request, u User) {
        if r.Method != http.MethodGet {
            http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
            return
        }

        renderApp(w, "static/app_home.html", AppBaseData{
            PageTitle: "Accueil",
            User:      u.Pseudo,
            Active:    "home",
        })
    })(w, r)
}