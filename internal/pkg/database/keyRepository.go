package database

import (
	"database/sql"
	"log"

	"github.com/phishhook/go-backend/internal/pkg/database/models"
)

type NoKeyExistsError struct {
}

func (e *NoKeyExistsError) Error() string {
	return "Invalid API key"
}

func (env *Env) GetApiKey(key string) (*models.Key, error) {
	var k models.Key
	// Gets a single row
	err := env.DB.QueryRow("SELECT key, last_used_at, active FROM api_keys WHERE key = $1", key).Scan(&k.Key, &k.LastUsedAt, &k.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Api key: %s does not exist", key)
			return nil, &NoKeyExistsError{}
		}
		log.Printf("Failed to scan rows from api_keys table: %s", err)
		return nil, err
	}
	return &k, nil
}
