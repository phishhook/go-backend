package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/phishhook/go-backend/internal/pkg/database/models"
)

type NoLinkExistsError struct {
	LinkIdentifier interface{}
}

func (e *NoLinkExistsError) Error() string {
	switch v := e.LinkIdentifier.(type) {
	case int:
		return fmt.Sprintf("Link with id %d does not exist", v)
	case string:
		return fmt.Sprintf("Link with URL %s does not exist", v)
	default:
		return "Invalid link identifier; must be an int ID or a string URL"
	}
}

// add a new link
func (env *Env) AddNewLink(userId int, url, isPhishing string) (int, error) {
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
	err := env.DB.QueryRow("INSERT INTO links (user_id, url, is_phishing) VALUES ($1, $2, $3) RETURNING id", userId, url, isPhishing).Scan(&link.ID)
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

// fetch all links for a specific user id
func (env *Env) GetLinksByUserId(userId string) ([]models.Link, error) {
	rows, err := env.DB.Query("SELECT id, user_id, url, clicked_at, is_phishing FROM links WHERE user_id = $1", userId)
	if err != nil {
		log.Println("Failed gather rows from table:", err)
		return nil, err
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		err := rows.Scan(&l.ID, &l.UserId, &l.Url, &l.ClickedAt, &l.IsPhishing)
		if err != nil {
			log.Println("Failed to scan link:", err)
			continue // or return depending on how you want to handle errors
		}
		links = append(links, l)
	}
	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		log.Println("Failed to iterate:", err)
		return nil, err
	}

	return links, nil
}

// fetch all links for a specific link id.
func (env *Env) GetLinkByLinkId(id string) (*models.Link, error) {
	var link models.Link
	// Gets a single row
	err := env.DB.QueryRow("SELECT id, user_id, url, clicked_at, is_phishing FROM links WHERE id = $1", id).Scan(&link.ID, &link.UserId, &link.Url, &link.ClickedAt, &link.IsPhishing)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No link found with id: %d", link.ID)
			return nil, &NoLinkExistsError{LinkIdentifier: id}
		}
		log.Printf("Failed to get the link with id: %d. Error: %s", link.ID, err)
		return nil, err
	}
	return &link, nil
}

func (env *Env) GetLinkByUrl(url string) (*models.Link, error) {
	var link models.Link
	// Gets a single row
	err := env.DB.QueryRow("SELECT id, user_id, url, clicked_at, is_phishing FROM links WHERE url = $1", url).Scan(&link.ID, &link.UserId, &link.Url, &link.ClickedAt, &link.IsPhishing)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No link found with url: %s", link.Url)
			return nil, &NoLinkExistsError{LinkIdentifier: url}
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
