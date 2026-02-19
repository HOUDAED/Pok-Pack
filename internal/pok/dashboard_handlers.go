package pok

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	RequireAuth(func(w http.ResponseWriter, r *http.Request, u User) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		collection, err := ListCardsByUser(u.ID)
		if err != nil {
			log.Printf("Erreur list cards: %v", err)
			http.Error(w, "Erreur serveur.", http.StatusInternalServerError)
			return
		}
		total, err := CountCardsByUser(u.ID)
		if err != nil {
			log.Printf("Erreur count cards: %v", err)
			http.Error(w, "Erreur serveur.", http.StatusInternalServerError)
			return
		}

		typeCounts := ComputeTypeCounts(collection)
		distinct := DistinctTypeCount(collection)
		chartJSON := TypeCountsToJSON(typeCounts)

		rarestName := ""
		rarestCount := 0
		if total > 0 {
			if rp, err := RarestPokemon(u.ID); err == nil {
				rarestName = rp.Name
				rarestCount = rp.Count
			}
		}

		renderApp(w, "web/templates/app_home.html", AppBaseData{
			PageTitle:      "Accueil",
			User:           u.Pseudo,
			Active:         "home",
			TotalCards:     total,
			DistinctTypes:  distinct,
			TypesChartJSON: chartJSON,
			RarestName:     rarestName,
			RarestCount:    rarestCount,
		})
	})(w, r)
}

func OpenPackHandler(w http.ResponseWriter, r *http.Request) {
	RequireAuth(func(w http.ResponseWriter, r *http.Request, u User) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/pack", http.StatusSeeOther)
			return
		}

		
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		opened := make([]Card, 0, 3)
		for len(opened) < 3 {
			id := 1 + rnd.Intn(1025)
			card, err := fetchPokemonCard(id)
			if err != nil {
				log.Printf("Erreur PokeAPI id=%d: %v", id, err)
				continue
			}
			opened = append(opened, card)
			if err := InsertCard(u.ID, card.PokemonID, card.Name, card.Types, card.ImageURL, time.Now().Unix()); err != nil {
				log.Printf("Erreur insertion card: %v", err)
				http.Error(w, "Erreur serveur.", http.StatusInternalServerError)
				return
			}
		}

		renderApp(w, "web/templates/app_pack.html", AppBaseData{
			PageTitle: "Ouvrir un pack",
			User:      u.Pseudo,
			Active:    "pack",
			Opened:    opened,
		})
	})(w, r)
}
