package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/phishhook/go-backend/internal/pkg/database/models"
)

type NoLinkExistsError struct {
	LinkId int
}

func (e *NoLinkExistsError) Error() string {
	return fmt.Sprintf("Link with id %d does not exist", e.LinkId)
}

// add a new link
func (env *Env) AddNewLink(url, isPhishing string) (int, error) {
	// Define valid statuses
	validStatuses := map[string]bool{
		"phishing":      true,
		"safe":          true,
		"indeterminate": true,
	}

	// Check if isPhishing is valid
	if _, exists := validStatuses[isPhishing]; !exists {
		return 0, fmt.Errorf("isPhishing must be 'phishing', 'safe', or 'indeterminate'")
	}

	var link models.Link
	err := env.DB.QueryRow("INSERT INTO links (url, is_phishing) VALUES ($1, $2) RETURNING id", url, isPhishing).Scan(&link.ID)
	if err != nil {
		log.Printf("Failed to add link: %s. Error: %s", url, err)
		return 0, err
	}

	return link.ID, nil
}

// get all links we have
func (env *Env) GetAllLinks() ([]models.Link, error) {
	rows, err := env.DB.Query("SELECT id, user_id, url, clicked_at, is_phishing FROM links")
	if err != nil {
		log.Println("Failed gather rows from users table:", err)
		return nil, err
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var link models.Link
		err := rows.Scan(&link.ID, &link.UserId, &link.Url, &link.ClickedAt, &link.IsPhishing)
		if err != nil {
			log.Println("Failed to scan link:", err)
			continue // or return depending on how you want to handle errors
		}
		links = append(links, link)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		log.Println("Failed to iterate over link rows:", err)
		return nil, err
	}

	return links, nil
}

// fetch all links for a user.
func (env *Env) GetLinkById(id string) (*models.Link, error) {
	var link models.Link
	// Gets a single row
	err := env.DB.QueryRow("SELECT id, user_id, url, clicked_at, is_phishing FROM links WHERE id = $1", id).Scan(&link.ID, &link.UserId, &link.Url, &link.ClickedAt, &link.IsPhishing)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No link found with id: %d", link.ID)
			return nil, &NoLinkExistsError{LinkId: link.ID}
		}
		log.Printf("Failed to get the link with id: %d. Error: %s", link.ID, err)
		return nil, err
	}
	return &link, nil
}

// delete a link
func (env *Env) DeleteLink(id string) (int, error) {
	return 0, nil
}
