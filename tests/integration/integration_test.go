package integration_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func getDatabaseConnection(t *testing.T) *sql.DB {
	t.Helper()
	// Generlize this so that it can be used to connect to any Postgres
	// database to run tests.
	connStr := "host=localhost user=dbadmin dbname=compliance password=secret port=5432 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		msg := fmt.Sprintf("Unable to initialize connection to test database: %s", err)
		t.Skip(msg)
	}

	// Wait up to 30 seconds to establish a connection with the database.
	// Remove this logic when we have the ability to set retries in the
	// database connection directly
	// (https://github.com/golang/go/issues/48309).
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err != nil {
			// We should only retry if we're dealing with a network
			// issue of some kind. No amount of retries is going to
			// fix incorrect credentials.
			var netError *net.OpError
			if errors.As(err, &netError) {
				t.Logf("Retrying database connection due to error: %s", err)
				time.Sleep(3 * time.Second)
				continue
			} else {
				msg := fmt.Sprintf("Unable to establish connection to test database: %s", err)
				t.Skip(msg)
			}
		}
	}

	return db
}

func getMigrationHelper(t *testing.T) *migrate.Migrate {
	t.Helper()
	db := getDatabaseConnection(t)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Skip("Unable to initialize database driver for migrations")
	}
	m, err := migrate.NewWithDatabaseInstance("file://../../migrations", "postgres", driver)
	if err != nil {
		t.Skip("Unable to initialize migrations")
	}
	return m
}

// Each test assumes the database is unmanaged. The test is responsible for
// setting up the state it requires for its test logic. This keeps the
// getMigrationHelper() method clean of any assumptions about what the tests
// expect from it.
func TestSubject(t *testing.T) { // nolint:paralleltest // database tests should not run in parallel
	getMigrationHelper(t)
	// m.Up()
	// Test logic
	// m.Down()
}

func TestMigration(t *testing.T) { // nolint:paralleltest // database tests should not run in parallel
	m := getMigrationHelper(t)

	version, dirty, err := m.Version()

	expectedVersion := uint(0)
	expectedState := false
	assert.Equal(t, expectedVersion, version, "Database version mismatch: want %d but got %d", expectedVersion, version)
	assert.Equal(t, expectedState, dirty, "Database state mismatch: want %t but got %t", expectedState, dirty)
	// Currently, Version() doesn't return a typed error, but a generic one
	// with a specific string (e.g., "no migration"). If this changes, we
	// should update the test to check the error type returned and that
	// it's what we expect.
	assert.NotEmpty(t, err, "Collecting the version should return an error on an empty database")

	// Upgrade the database and make sure all upgrades apply cleanly.
	err = m.Up()
	version, dirty, _ = m.Version()
	expectedVersion = uint(1)
	assert.Equal(t, expectedVersion, version, "Database version mismatch: want %d but got %d", expectedVersion, version)
	assert.Equal(t, false, dirty, "Database state mismatch: want %t but got %t", false, dirty)
	assert.Equal(t, err, nil, "Error upgrading the database: %s", err)

	// Downgrade the database back to 0 and make sure all downgrades apply cleanly.
	err = m.Down()
	version, dirty, _ = m.Version()
	expectedVersion = uint(0)
	assert.Equal(t, expectedVersion, version, "Database version mismatch: want %d but got %d", expectedVersion, version)
	assert.Equal(t, false, dirty, "Database state mismatch: want %t but got %t", false, dirty)
	assert.Equal(t, err, nil, "Error downgrading the database: %s", err)
}
