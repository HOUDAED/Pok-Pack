package pok

import (
	"encoding/json"
	"sort"
	"strings"
)

type TypeCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// pour ameliorer le rendu du project j'ai ajouté des fonctions pour calculer les stats de types à partir des cartes obtenues par l'utilisateur. 

func splitTypes(types string) []string {
	parts := strings.Split(types, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.ToLower(strings.TrimSpace(p))
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}

func ComputeTypeCounts(cards []Card) map[string]int {
	m := map[string]int{}
	for _, c := range cards {
		for _, t := range splitTypes(c.Types) {
			m[t]++
		}
	}
	return m
}

func DistinctTypeCount(cards []Card) int {
	seen := map[string]struct{}{}
	for _, c := range cards {
		for _, t := range splitTypes(c.Types) {
			seen[t] = struct{}{}
		}
	}
	return len(seen)
}

func TypeCountsToJSON(m map[string]int) string {
	items := make([]TypeCount, 0, len(m))
	for k, v := range m {
		items = append(items, TypeCount{Type: k, Count: v})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Type < items[j].Type
		}
		return items[i].Count > items[j].Count
	})
	
	if len(items) > 10 {
		items = items[:10]
	}
	b, _ := json.Marshal(items)
	return string(b)
}
