package pok

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func PackPageHandler(w http.ResponseWriter, r *http.Request) {
	RequireAuth(func(w http.ResponseWriter, r *http.Request, u User) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/pack", http.StatusSeeOther)
			return
		}
		renderApp(w, "static/app_pack.html", AppBaseData{
			PageTitle: "Ouvrir un pack",
			User:      u.Pseudo,
			Active:    "pack",
			Opened:    []Card{},
		})
	})(w, r)
}

func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	RequireAuth(func(w http.ResponseWriter, r *http.Request, u User) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/collection", http.StatusSeeOther)
			return
		}
		collection, err := ListCardsByUser(u.ID)
		if err != nil {
			log.Printf("Erreur list cards: %v", err)
			http.Error(w, "Erreur serveur.", http.StatusInternalServerError)
			return
		}

		notice := ""
		if r.URL.Query().Get("added") == "1" {
			notice = "Carte ajoutée à ta collection."
		}
		renderApp(w, "static/app_collection.html", AppBaseData{
			PageTitle:     "Collection",
			User:          u.Pseudo,
			Active:        "collection",
			NoticeSuccess: notice,
			Collection:    collection,
		})
	})(w, r)
}

func AddCardHandler(w http.ResponseWriter, r *http.Request) {
	RequireAuth(func(w http.ResponseWriter, r *http.Request, u User) {
		switch r.Method {
		case http.MethodGet:
			renderApp(w, "static/app_add.html", AppBaseData{
				PageTitle: "Ajouter une carte",
				User:      u.Pseudo,
				Active:    "add",
			})
			return
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				renderApp(w, "static/app_add.html", AppBaseData{
					PageTitle:     "Ajouter une carte",
					User:          u.Pseudo,
					Active:        "add",
					NoticeError:   "Formulaire invalide.",
					NoticeSuccess: "",
				})
				return
			}
			idStr := r.FormValue("pokemon_id")
			pid, err := strconv.Atoi(idStr)
			if err != nil || pid < 1 || pid > 1025 {
				renderApp(w, "static/app_add.html", AppBaseData{
					PageTitle:   "Ajouter une carte",
					User:        u.Pseudo,
					Active:      "add",
					NoticeError: "ID Pokémon invalide (1 à 1025).",
				})
				return
			}

			card, err := fetchPokemonCard(pid)
			if err != nil {
				log.Printf("Erreur PokeAPI add id=%d: %v", pid, err)
				renderApp(w, "static/app_add.html", AppBaseData{
					PageTitle:   "Ajouter une carte",
					User:        u.Pseudo,
					Active:      "add",
					NoticeError: "Impossible de récupérer ce Pokémon via PokeAPI.",
				})
				return
			}
			if err := InsertCard(u.ID, card.PokemonID, card.Name, card.Types, card.ImageURL, time.Now().Unix()); err != nil {
				log.Printf("Erreur insertion card: %v", err)
				http.Error(w, "Erreur serveur.", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/collection?added=1", http.StatusSeeOther)
			return
		default:
			http.Redirect(w, r, "/card/add", http.StatusSeeOther)
		}
	})(w, r)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	RequireAuth(func(w http.ResponseWriter, r *http.Request, u User) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/stats", http.StatusSeeOther)
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

		mostName := ""
		mostCount := 0
		mostImg := ""
		mostTypes := ""
		if total > 0 {
			if mp, err := MostObtainedPokemon(u.ID); err == nil {
				mostName = mp.Name
				mostCount = mp.Count
				mostImg = mp.ImageURL
				mostTypes = mp.Types
			}
		}

		renderApp(w, "static/app_stats.html", AppBaseData{
			PageTitle:         "Statistiques",
			User:              u.Pseudo,
			Active:            "stats",
			TotalCards:        total,
			DistinctTypes:     distinct,
			TypesChartJSON:    chartJSON,
			MostObtainedName:  mostName,
			MostObtainedCount: mostCount,
			MostObtainedImage: mostImg,
			MostObtainedTypes: mostTypes,
		})
	})(w, r)
}
