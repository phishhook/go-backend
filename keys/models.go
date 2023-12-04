package keys

import (
	"database/sql"
	"log"
	"time"
)

type Key struct {
	ID         int       `json:"id,omitempty"`
	Key        string    `json:"key"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
	Active     bool      `json:"active"`
}

func ApiKey(db *sql.DB, key string) (*Key, error) {
	var k Key
	// Gets a single row
	err := db.QueryRow("SELECT key, last_used_at, active FROM api_keys WHERE key = $1", key).Scan(&k.Key, &k.LastUsedAt, &k.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Api key: %s does not exist", key)
			return nil, err
		}
		log.Printf("Failed to scan rows from api_keys table: %s", err)
		return nil, err
	}
	return &k, nil
}
