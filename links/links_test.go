package links

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

const (
	EXPECTED_URL            = "http://example.com"
	ErrUnfilledExpectations = "Error was not expected while fetching links: %s"
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

func TestLinksIndex(t *testing.T) {
	// Create a new mock database
	db, mock := CreateMockDB(t)
	defer db.Close()

	// Set up mock rows
	rows := sqlmock.NewRows([]string{"id", "user_id", "url", "clicked_at", "is_phishing", "percentage"}).
		AddRow(1, 100, EXPECTED_URL, FIXED_TIME, "safe", "50.00")

	// Expectations
	mock.ExpectQuery("^SELECT (.+) FROM links$").WillReturnRows(rows)

	// Call AllLinks
	links, err := AllLinks(db)
	if err != nil {
		t.Errorf("error was not expected while fetching links: %s", err)
	}

	// Assertions
	expectedLinks := []*Link{
		{ID: 1, UserId: 100, Url: EXPECTED_URL, ClickedAt: FIXED_TIME, IsPhishing: "safe", Percentage: "50.00"},
	}
	if !reflect.DeepEqual(links, expectedLinks) {
		t.Errorf("expected links %v, but got %v", expectedLinks, links)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(ErrUnfilledExpectations, err)
	}
}

func TestPostLink(t *testing.T) {
	// Create a new mock database
	db, mock := CreateMockDB(t)
	defer db.Close()

	mock.ExpectQuery("^INSERT INTO links (.+) VALUES (.+)$").
		WithArgs(1, EXPECTED_URL, "safe", "50.00").               // Ensure these match the arguments in AddNewLink
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Simulate RETURNING id

	// Call AddNewLink
	id, err := AddNewLink(db, 1, EXPECTED_URL, "safe", "50.00")
	if err != nil {
		t.Errorf("error was not expected while adding a new link: %s", err)
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

func TestGetLinksByUserID(t *testing.T) {
	// Create a new mock database
	db, mock := CreateMockDB(t)
	defer db.Close()

	// Set up mock rows
	rows := sqlmock.NewRows([]string{"id", "user_id", "url", "clicked_at", "is_phishing", "percentage"}).
		AddRow(1, 100, EXPECTED_URL, FIXED_TIME, "safe", "50.00")

	// Expectations
	mock.ExpectQuery("^SELECT (.+) FROM links WHERE user_id = \\$1$").
		WithArgs("1"). // Expect a string "1" instead of an integer 1
		WillReturnRows(rows)

	// Call LinksByUserId
	links, err := LinksByUserId(db, "1")
	if err != nil {
		t.Errorf("error was not expected while fetching links: %s", err)
	}

	// Assertions
	expectedLinks := []*Link{
		{ID: 1, UserId: 100, Url: EXPECTED_URL, ClickedAt: FIXED_TIME, IsPhishing: "safe", Percentage: "50.00"},
	}
	if !reflect.DeepEqual(links, expectedLinks) {
		t.Errorf("expected links %v, but got %v", expectedLinks, links)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(ErrUnfilledExpectations, err)
	}
}

func TestGetLinkByUrl(t *testing.T) {
	// Create a new mock database
	db, mock := CreateMockDB(t)
	defer db.Close()

	// Set up mock rows
	rows := sqlmock.NewRows([]string{"id", "user_id", "url", "clicked_at", "is_phishing", "percentage"}).
		AddRow(1, 100, EXPECTED_URL, FIXED_TIME, "safe", "50.00")

	// Expectations
	mock.ExpectQuery("^SELECT (.+) FROM links WHERE url = \\$1$").
		WithArgs(EXPECTED_URL). // Expect a string "1" instead of an integer 1
		WillReturnRows(rows)

	// Call LinksByUserId
	link, err := LinkByUrl(db, EXPECTED_URL)
	if err != nil {
		t.Errorf("error was not expected while fetching links: %s", err)
	}

	// Assertions
	expectedLinks := &Link{
		ID: 1, UserId: 100, Url: EXPECTED_URL, ClickedAt: FIXED_TIME, IsPhishing: "safe", Percentage: "50.00",
	}
	if !reflect.DeepEqual(link, expectedLinks) {
		t.Errorf("expected link %v, but got %v", expectedLinks, link)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(ErrUnfilledExpectations, err)
	}
}
