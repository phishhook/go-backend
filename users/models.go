package users

import (
	"database/sql"
	"log"
	"time"
)

type User struct {
	ID             int       `json:"id,omitempty"`
	Username       string    `json:"username"`
	PhoneNumber    string    `json:"phone_number"`
	CreatedAt      time.Time `json:"created_at"`
	AnonymizeLinks bool      `json:"anonymize_links"`
}

// Fetching all user profiles that we have
func AllUsers(db *sql.DB) ([]*User, error) {
	rows, err := db.Query("SELECT id, username, phone_number, created_at, anonymize_links FROM users")
	if err != nil {
		log.Println("Failed gather rows from users table:", err)
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		u := new(User)
		err := rows.Scan(&u.ID, &u.Username, &u.PhoneNumber, &u.CreatedAt, &u.AnonymizeLinks)
		if err != nil {
			log.Println("Failed to scan user:", err)
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// Add a new user.
func AddNewUser(db *sql.DB, phoneNumber, username string, anonLinks bool) (int, error) {
	var user User
	// Check if the user already exists.
	isUserAlreadyCreated, err := IsUserAlreadyCreated(db, phoneNumber)
	if err != nil {
		log.Println("Failed to check if user already exists:", err)
		return 0, err
	}
	if isUserAlreadyCreated {
		return 0, &UserAlreadyExistsError{PhoneNumber: phoneNumber}
	}

	err = db.QueryRow("INSERT INTO users (username, phone_number, anonymize_links) VALUES ($1, $2, $3) RETURNING id", username, phoneNumber, anonLinks).Scan(&user.ID)
	if err != nil {
		log.Printf("Failed to add user: %s. Error: %s", username, err)
		return 0, err
	}

	return user.ID, nil
}

// get a specific user by their phone number.
func UserByPhoneNumber(db *sql.DB, phoneNumber string) (*User, error) {
	var u User
	// Gets a single row
	err := db.QueryRow("SELECT id, username, phone_number, created_at, anonymize_links FROM users WHERE phone_number = $1", phoneNumber).Scan(&u.ID, &u.Username, &u.PhoneNumber, &u.CreatedAt, &u.AnonymizeLinks)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with phone_number: %s", phoneNumber)
			return nil, err
		}
		log.Printf("Failed to get the user with phone_number: %s. Error: %s", phoneNumber, err)
		return nil, err
	}
	return &u, nil
}

// Validating user credentials.
func IsUserAlreadyCreated(db *sql.DB, phoneNumber string) (bool, error) {
	_, err := UserByPhoneNumber(db, phoneNumber)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User does not exist. Adding them to the database.")
			return false, nil // our one success case, we can create the user then.
		}
		log.Println("Failed to validate user:", err)
		return false, err // Some other error occurred.
	}

	return true, nil // user already exists.
}
