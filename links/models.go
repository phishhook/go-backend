package links

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"
)

type Link struct {
	ID         int       `json:"link_id,omitempty"`
	UserId     int       `json:"user_id"`
	Url        string    `json:"url"`
	UrlScheme  string    `json:"url_scheme"`
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
func AddNewLink(db *sql.DB, userId int, inputURL, isPhishing, percentage string) (int, error) {
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

	// parse out the http or https
	parsedUrl, err := url.Parse(inputURL)
	if err != nil {
		return 0, err // return an empty string and the error
	}
	urlScheme := parsedUrl.Scheme

	var link Link
	err = db.QueryRow("INSERT INTO links (user_id, url, is_phishing, percentage, url_scheme) VALUES ($1, $2, $3, $4, %5) RETURNING id", userId, inputURL, isPhishing, percentage, urlScheme).Scan(&link.ID)
	if err != nil {
		log.Printf("Failed to add link: %s. Error: %s", inputURL, err)
		return 0, err
	}

	return link.ID, nil
}

// fetch all links for a specific user id
func LinksByUserId(db *sql.DB, userId string) ([]*Link, error) {
	rows, err := db.Query("SELECT id, user_id, url, clicked_at, is_phishing, percentage, url_scheme FROM links WHERE user_id = $1", userId)
	if err != nil {
		log.Println("Failed gather rows from users table:", err)
		return nil, err
	}
	defer rows.Close()

	links := make([]*Link, 0)
	for rows.Next() {
		l := new(Link)
		err := rows.Scan(&l.ID, &l.UserId, &l.Url, &l.ClickedAt, &l.IsPhishing, &l.Percentage, &l.UrlScheme)
		if err != nil {
			log.Println("Failed to scan link:", err)
			return nil, err
		}
		links = append(links, l)
	}

	return links, nil
}

// fetch the link that corresponds to the given link url
func LinkByUrl(db *sql.DB, inputURL string) (*Link, error) {
	link := new(Link)

	// parse out the http or https
	parsedUrl, err := url.Parse(inputURL)
	if err != nil {
		log.Println("Failed to parse link:", err)
		return nil, err
	}
	urlScheme := parsedUrl.Scheme
	urlWithoutScheme := parsedUrl.Host + parsedUrl.Path

	err = db.QueryRow("SELECT id, user_id, url, url_scheme, clicked_at, is_phishing, percentage FROM links WHERE url = $1 AND url_scheme = $2", urlWithoutScheme, urlScheme).Scan(&link.ID, &link.UserId, &link.Url, &link.UrlScheme, &link.ClickedAt, &link.IsPhishing, &link.Percentage)
	if err != nil {
		log.Printf("Failed to get the link with url: %s and urlScheme: %s. Error: %s", urlWithoutScheme, urlScheme, err)
		return nil, err
	}
	return link, nil
}
