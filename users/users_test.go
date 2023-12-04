package users

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

const (
	USERNAME                = "gandalf"
	PHONE_NUMBER            = "16512529620"
	ErrUnfilledExpectations = "Error was not expected while fetching links: %s"
	ErrFetchingLinks        = "error was not expected while fetching links: %s"
	ErrUnmetExpectations    = "unmet expectation error: %s"
	InsertUserSqlStatement  = "^SELECT id, username, phone_number, created_at, anonymize_links FROM users WHERE phone_number = \\$1"
)

var FIXED_TIME = time.Now()

// CreateMockDB sets up a mock database and returns it along with the sqlmock
func CreateMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestUsersIndex(t *testing.T) {
	// Create a new mock database
	db, mock := CreateMockDB(t)
	defer db.Close()

	// Set up mock rows
	rows := sqlmock.NewRows([]string{"id", "username", "phone_number", "created_at", "anonymize_links"}).
		AddRow(1, USERNAME, PHONE_NUMBER, FIXED_TIME, false)

	// Expectations
	mock.ExpectQuery("^SELECT (.+) FROM users$").WillReturnRows(rows)

	// Call AllLinks
	links, err := AllUsers(db)
	if err != nil {
		t.Errorf(ErrFetchingLinks, err)
	}

	// Assertions
	expectedLinks := []*User{
		{ID: 1, Username: USERNAME, PhoneNumber: PHONE_NUMBER, CreatedAt: FIXED_TIME, AnonymizeLinks: false},
	}
	if !reflect.DeepEqual(links, expectedLinks) {
		t.Errorf("expected links %v, but got %v", expectedLinks, links)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(ErrUnfilledExpectations, err)
	}
}

func TestAddNewUserNoUserFound(t *testing.T) {
	// Create a new mock database
	db, mock := CreateMockDB(t)
	defer db.Close()

	// Add mock expectation for the SELECT query
	mock.ExpectQuery(InsertUserSqlStatement).
		WithArgs(PHONE_NUMBER).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "phone_number", "created_at", "anonymize_links"})) // Simulate no existing user

	mock.ExpectQuery("^INSERT INTO users (.+) VALUES (.+)$").
		WithArgs(USERNAME, PHONE_NUMBER, false).                  // Ensure these match the arguments in AddNewUser
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Simulate RETURNING id

	// Call AddNewUser
	id, err := AddNewUser(db, PHONE_NUMBER, USERNAME, false)
	if err != nil {
		t.Errorf("error was not expected while adding a new user: %s", err)
	}

	// Assertions
	expectedID := 1
	if id != expectedID {
		t.Errorf("expected id %v, but got %v", expectedID, id)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(ErrUnfilledExpectations, err)
	}
}

func TestAddNewUserUserFound(t *testing.T) {
	// Create a new mock database
	db, mock := CreateMockDB(t)
	defer db.Close()

	// Simulate an existing user
	mock.ExpectQuery(InsertUserSqlStatement).
		WithArgs(PHONE_NUMBER).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "phone_number", "created_at", "anonymize_links"}).
			AddRow(1, "existingUser", PHONE_NUMBER, time.Now(), false))

	// Call AddNewUser
	_, err := AddNewUser(db, PHONE_NUMBER, USERNAME, false)

	// Assert that an error is returned
	if err == nil {
		t.Errorf("expected an error when adding a new user with an existing phone number, but got none")
	}

	// Optionally, assert that the error is of the specific type UserAlreadyExistsError
	if _, ok := err.(*UserAlreadyExistsError); !ok {
		t.Errorf("expected a UserAlreadyExistsError, but got a different type of error: %s", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(ErrUnmetExpectations, err)
	}
}

func TestUserByPhoneNumberFound(t *testing.T) {
	// Create a new mock database
	db, mock := CreateMockDB(t)
	defer db.Close()

	// Simulate a user being found
	mock.ExpectQuery(InsertUserSqlStatement).
		WithArgs(PHONE_NUMBER).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "phone_number", "created_at", "anonymize_links"}).
			AddRow(1, "testUser", PHONE_NUMBER, time.Now(), false))

	// Call UserByPhoneNumber
	user, err := UserByPhoneNumber(db, PHONE_NUMBER)
	if err != nil {
		t.Errorf("error was not expected while retrieving a user: %s", err)
	}
	if user == nil {
		t.Errorf("expected a user to be returned, but got nil")
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(ErrUnmetExpectations, err)
	}
}

func TestUserByPhoneNumberNotFound(t *testing.T) {
	// Create a new mock database
	db, mock := CreateMockDB(t)
	defer db.Close()

	// Simulate no user being found
	mock.ExpectQuery(InsertUserSqlStatement).
		WithArgs(PHONE_NUMBER).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "phone_number", "created_at", "anonymize_links"}))

	// Call UserByPhoneNumber
	user, err := UserByPhoneNumber(db, PHONE_NUMBER)
	if err == nil {
		t.Errorf("expected an error when no user is found, but got none")
	}
	if user != nil {
		t.Errorf("expected no user to be returned, but got a user")
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(ErrUnmetExpectations, err)
	}
}
