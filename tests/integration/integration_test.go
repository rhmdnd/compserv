package integration_test

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func setUpDatabase(t *testing.T) *sql.DB {
	t.Helper()
	// Generlize this so that it can be used to connect to any Postgres
	// database to run tests.
	connStr := "host=localhost user=dbadmin dbname=compliance password=secret port=5432 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		msg := fmt.Sprintf("Unable to initialize connection to test database: %s", err)
		t.Skip(msg)
	}
	if err := db.Ping(); err != nil {
		msg := fmt.Sprintf("Unable to establish connection to test database: %s", err)
		t.Skip(msg)
	}

	// Finish adding migration logic.
	return db
}

func TestSubject(t *testing.T) {
	setUpDatabase(t)
	t.Parallel()
}
