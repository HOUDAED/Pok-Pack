package pok

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	queryCacheMu sync.Mutex
	queryCache   = map[string]string{}
)

func LoadSQL(relPath string) (string, error) {
	clean := filepath.Clean(relPath)
	clean = strings.TrimPrefix(clean, string(filepath.Separator))

	queryCacheMu.Lock()
	if v, ok := queryCache[clean]; ok {
		queryCacheMu.Unlock()
		return v, nil
	}
	queryCacheMu.Unlock()

	b, err := os.ReadFile(clean)
	if err != nil {
		return "", fmt.Errorf("read sql %s: %w", clean, err)
	}

	sqlText := strings.TrimSpace(string(b))
	if sqlText == "" {
		return "", fmt.Errorf("sql file is empty: %s", clean)
	}

	queryCacheMu.Lock()
	queryCache[clean] = sqlText
	queryCacheMu.Unlock()

	return sqlText, nil
}
