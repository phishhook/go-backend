package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/phishhook/go-backend/internal/pkg/database/models"
)

type UserAlreadyExistsError struct {
	PhoneNumber string
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user with phone number %s already exists", e.PhoneNumber)
}

type NoUserExistsError struct {
	PhoneNumber string
}

func (e *NoUserExistsError) Error() string {
	return fmt.Sprintf("user with phone number %s does not exist", e.PhoneNumber)
}

// Creating a new user.
func (env *Env) AddNewUser(phoneNumber, username string) (int, error) {
	var user models.User
	// Check if the user already exists.
	isUserAlreadyCreated, err := env.IsUserAlreadyCreated(phoneNumber)
	if err != nil {
		log.Println("Failed to check if user already exists:", err)
		return 0, err
	}
	if isUserAlreadyCreated {
		return 0, &UserAlreadyExistsError{PhoneNumber: phoneNumber}
	}

	err = env.DB.QueryRow("INSERT INTO users (username, phone_number) VALUES ($1, $2) RETURNING id", username, phoneNumber).Scan(&user.ID)
	if err != nil {
		log.Printf("Failed to add user: %s. Error: %s", username, err)
		return 0, err
	}

	return user.ID, nil
}

// Fetching all user profiles.
func (env *Env) GetAllUsers() ([]models.User, error) {
	rows, err := env.DB.Query("SELECT id, username, phone_number, created_at FROM users")
	if err != nil {
		log.Println("Failed gather rows from users table:", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Username, &u.PhoneNumber, &u.CreatedAt)
		if err != nil {
			log.Println("Failed to scan user:", err)
			continue // or return depending on how you want to handle errors
		}
		users = append(users, u)
	}
	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		log.Println("Failed to iterate over user rows:", err)
		return nil, err
	}

	return users, nil
}

// get a specific user by their phone number.
func (env *Env) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	var u models.User
	// Gets a single row
	err := env.DB.QueryRow("SELECT id, username, phone_number, created_at FROM users WHERE phone_number = $1", phoneNumber).Scan(&u.ID, &u.Username, &u.PhoneNumber, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with phone_number: %s", phoneNumber)
			return nil, &NoUserExistsError{PhoneNumber: phoneNumber}
		}
		log.Printf("Failed to get the user with phone_number: %s. Error: %s", phoneNumber, err)
		return nil, err
	}
	return &u, nil
}

// Updating user information.
func (env *Env) UpdateUser() ([]models.User, error) {
	return nil, nil
}

// Validating user credentials.
func (env *Env) IsUserAlreadyCreated(phoneNumber string) (bool, error) {
	_, err := env.GetUserByPhoneNumber(phoneNumber)

	if err != nil {
		var NoUserExistsErr *NoUserExistsError
		if errors.As(err, &NoUserExistsErr) {
			log.Println("User does not exist.", err)
			return false, nil // our one success case, we can create the user then.
		}
		log.Println("Failed to validate user:", err)
		return false, err // Some other error occurred.
	}

	return true, nil // user already exists.
}
