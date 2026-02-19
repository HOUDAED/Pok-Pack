package pok

import (
	"database/sql"
)

type User struct {
	ID           int
	Pseudo       string
	Email        string
	PasswordHash string
}

func GetUserByPseudoOrEmail(login string) (User, error) {
	var u User
	err := pokdb.QueryRow(
		"SELECT id, pseudo, email, password_hash FROM users WHERE pseudo = ? OR email = ?",
		login,
		login,
	).Scan(&u.ID, &u.Pseudo, &u.Email, &u.PasswordHash)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func GetUserBySessionID(sessionID string, nowUnix int64) (User, error) {
	var u User
	err := pokdb.QueryRow(
		`SELECT u.id, u.pseudo, u.email, u.password_hash
		 FROM users u
		 JOIN sessions s ON s.user_id = u.id
		 WHERE s.id = ? AND s.expires_at > ?`,
		sessionID,
		nowUnix,
	).Scan(&u.ID, &u.Pseudo, &u.Email, &u.PasswordHash)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}
