package pok

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type pokeAPIResp struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Sprites struct {
		FrontDefault string         `json:"front_default"`
		Other        map[string]any `json:"other"`
	} `json:"sprites"`
}

func fetchPokemonCard(pokemonID int) (Card, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", pokemonID)
	client := &http.Client{Timeout: 8 * time.Second}

	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return Card{}, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "MyPokePack/1.0")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < 3 {
				time.Sleep(time.Duration(160*attempt) * time.Millisecond)
				continue
			}
			return Card{}, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			_ = resp.Body.Close()
			lastErr = fmt.Errorf("pokeapi status: %s", resp.Status)
			
			if (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500) && attempt < 3 {
				time.Sleep(time.Duration(240*attempt) * time.Millisecond)
				continue
			}
			return Card{}, lastErr
		}

		var p pokeAPIResp
		if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
			_ = resp.Body.Close()
			lastErr = err
			if attempt < 3 {
				time.Sleep(time.Duration(160*attempt) * time.Millisecond)
				continue
			}
			return Card{}, err
		}
		_ = resp.Body.Close()

		types := ""
		for i, t := range p.Types {
			if i > 0 {
				types += ", "
			}
			types += t.Type.Name
		}

		imageURL := p.Sprites.FrontDefault
		
		if p.Sprites.Other != nil {
			if oaAny, ok := p.Sprites.Other["official-artwork"]; ok {
				if m, ok := oaAny.(map[string]any); ok {
					if v, ok := m["front_default"].(string); ok && v != "" {
						imageURL = v
					}
				}
			}
		}

		return Card{PokemonID: p.ID, Name: p.Name, Types: types, ImageURL: imageURL}, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("pokeapi error")
	}
	return Card{}, lastErr
}
