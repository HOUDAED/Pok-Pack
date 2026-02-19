package pok

import (
	"database/sql"
)

func CountCardsByUser(userID int) (int, error) {
	var n int
	err := pokdb.QueryRow("SELECT COUNT(*) FROM cards WHERE user_id = ?", userID).Scan(&n)
	return n, err
}

func ListCardsByUser(userID int) ([]Card, error) {
	rows, err := pokdb.Query("SELECT pokemon_id, name, types, image_url FROM cards WHERE user_id = ? ORDER BY id DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Card
	for rows.Next() {
		var c Card
		if err := rows.Scan(&c.PokemonID, &c.Name, &c.Types, &c.ImageURL); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		out = []Card{}
	}
	return out, nil
}

func HasAnyCard(userID int) (bool, error) {
	var id int
	err := pokdb.QueryRow("SELECT id FROM cards WHERE user_id = ? LIMIT 1", userID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

type PokemonCount struct {
	Card
	Count int
}

func MostObtainedPokemon(userID int) (PokemonCount, error) {
	var pc PokemonCount
	err := pokdb.QueryRow(
		`SELECT pokemon_id, name, types, image_url, COUNT(*) as n
		 FROM cards
		 WHERE user_id = ?
		 GROUP BY pokemon_id
		 ORDER BY n DESC
		 LIMIT 1`,
		userID,
	).Scan(&pc.PokemonID, &pc.Name, &pc.Types, &pc.ImageURL, &pc.Count)
	if err != nil {
		return PokemonCount{}, err
	}
	return pc, nil
}

func RarestPokemon(userID int) (PokemonCount, error) {
	var pc PokemonCount
	err := pokdb.QueryRow(
		`SELECT pokemon_id, name, types, image_url, COUNT(*) as n
		 FROM cards
		 WHERE user_id = ?
		 GROUP BY pokemon_id
		 ORDER BY n ASC
		 LIMIT 1`,
		userID,
	).Scan(&pc.PokemonID, &pc.Name, &pc.Types, &pc.ImageURL, &pc.Count)
	if err != nil {
		return PokemonCount{}, err
	}
	return pc, nil
}
