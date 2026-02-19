package pok

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var pokdb *sql.DB

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de la base de données (%s) : %v\n", filepath, err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Printf("Erreur lors du test de connexion (Ping): %v\n", err)
		db.Close()
		log.Printf("Erreur ouverture DB (%s) : %v\n", filepath, err)

		return nil, err
	}

	log.Println("Connexion SQLite établie et vérifiée.")

	// Créer les tables
	for name, query := range TablesSQL {
		log.Printf("Création table '%s'...", name)
		if _, err = db.Exec(query); err != nil {
			log.Printf("Erreur création table '%s' : %v\n", name, err)
			db.Close()
			return nil, err
		}
	}
	pokdb = db
	log.Println("Base de données initialisée avec succès.")
	return db, nil
}
