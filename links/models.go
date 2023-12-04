package links

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Link struct {
	ID         int       `json:"link_id,omitempty"`
	UserId     int       `json:"user_id"`
	Url        string    `json:"url"`
	ClickedAt  time.Time `json:"clicked_at"`
	IsPhishing string    `json:"is_phishing"`
	Percentage string    `json:"percentage"`
}

// get all links we have
func AllLinks(db *sql.DB) ([]*Link, error) {
	rows, err := db.Query("SELECT id, user_id, url, clicked_at, is_phishing, percentage FROM links")
	if err != nil {
		log.Println("Failed gather rows from users table:", err)
		return nil, err
	}
	defer rows.Close()

	links := make([]*Link, 0)
	for rows.Next() {
		lk := new(Link)
		err := rows.Scan(&lk.ID, &lk.UserId, &lk.Url, &lk.ClickedAt, &lk.IsPhishing, &lk.Percentage)
		if err != nil {
			log.Println("Failed to scan link:", err)
			return nil, err
		}
		links = append(links, lk)
	}
	return links, nil
}

// add a new link
func AddNewLink(db *sql.DB, userId int, url, isPhishing, percentage string) (int, error) {
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

	// Check that percentage is a valid 5 char string
	if len(percentage) != 5 {
		return 0, fmt.Errorf("percentage must be a 5 char string. Example: '50.00'")
	}

	var link Link
	err := db.QueryRow("INSERT INTO links (user_id, url, is_phishing, percentage) VALUES ($1, $2, $3, $4) RETURNING id", userId, url, isPhishing, percentage).Scan(&link.ID)
	if err != nil {
		log.Printf("Failed to add link: %s. Error: %s", url, err)
		return 0, err
	}

	return link.ID, nil
}

// fetch all links for a specific user id
func LinksByUserId(db *sql.DB, userId string) ([]*Link, error) {
	rows, err := db.Query("SELECT id, user_id, url, clicked_at, is_phishing, percentage FROM links WHERE user_id = $1", userId)
	if err != nil {
		log.Println("Failed gather rows from users table:", err)
		return nil, err
	}
	defer rows.Close()

	links := make([]*Link, 0)
	for rows.Next() {
		l := new(Link)
		err := rows.Scan(&l.ID, &l.UserId, &l.Url, &l.ClickedAt, &l.IsPhishing, &l.Percentage)
		if err != nil {
			log.Println("Failed to scan link:", err)
			return nil, err
		}
		links = append(links, l)
	}

	return links, nil
}

// fetch the link that corresponds to the given link url
func LinkByUrl(db *sql.DB, url string) (*Link, error) {
	link := new(Link)
	err := db.QueryRow("SELECT id, user_id, url, clicked_at, is_phishing FROM links WHERE url = $1", url).Scan(&link.ID, &link.UserId, &link.Url, &link.ClickedAt, &link.IsPhishing, &link.Percentage)
	if err != nil {
		log.Printf("Failed to get the link with id: %d. Error: %s", link.ID, err)
		return nil, err
	}
	return link, nil
}
