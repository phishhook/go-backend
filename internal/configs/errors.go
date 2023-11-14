package errors

import "fmt"

// userRepo Errors
type UserAlreadyExistsError struct {
	PhoneNumber string
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user with phone number %s already exists", e.PhoneNumber)
}
