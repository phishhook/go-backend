package users

import "fmt"

type UserAlreadyExistsError struct {
	PhoneNumber string
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user with phone number %s already exists", e.PhoneNumber)
}
