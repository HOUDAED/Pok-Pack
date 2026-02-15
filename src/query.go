package pok

// TablesSQL contient toutes les requêtes de création de tables
var TablesSQL = map[string]string{
	"users": `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pseudo TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);`,
	"sessions": `CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires_at INTEGER NOT NULL,
		created_at INTEGER NOT NULL
	);`,
	"cards": `CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		pokemon_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		types TEXT NOT NULL,
		image_url TEXT NOT NULL,
		obtained_at INTEGER NOT NULL
	);`,
}

func InsertValuesUser(pseudo, email, passwordHash string) error {
	_, err := pokdb.Exec("INSERT INTO users (pseudo, email, password_hash) VALUES (?, ?, ?)", pseudo, email, passwordHash)
	return err
}

func InsertSession(id string, userID int, expiresAt int64, createdAt int64) error {
	_, err := pokdb.Exec("INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)", id, userID, expiresAt, createdAt)
	return err
}

func DeleteSession(id string) error {
	_, err := pokdb.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

func DeleteAllSessions() error {
	_, err := pokdb.Exec("DELETE FROM sessions")
	return err
}

