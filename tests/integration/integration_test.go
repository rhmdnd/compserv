package integration_test // nolint:typecheck

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Necessary to invoke migrations locally in the repository
	"github.com/google/uuid"
	_ "github.com/lib/pq" // Necessary to use the PostgreSQL database driver and connecting to a PostgreSQL database
	"github.com/stretchr/testify/assert"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const clusterName = "cluster.example.com"

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
				// Linting says we shouldn't use the following:
				// time.Sleep(3 * time.Second)
				// but we can't use
				// duration := 3
				// time.Sleep(duration * time.Second)
				// which causes a type mismatch.
				duration, _ := time.ParseDuration("0m3s")
				time.Sleep(duration)
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

func getGormHelper() *gorm.DB {
	connStr := "host=localhost user=dbadmin dbname=compliance password=secret port=5432 sslmode=disable"
	gormDB, _ := gorm.Open(gorm_postgres.Open(connStr), &gorm.Config{})
	return gormDB
}

func getUUIDString() string {
	value, _ := uuid.NewRandom()
	return value.String()
}

// Each test assumes the database is unmanaged. The test is responsible for
// setting up the state it requires for its test logic. This keeps the
// getMigrationHelper() method clean of any assumptions about what the tests
// expect from it. Each test must be run serially they have the power to change
// database schema and affect other tests.
func TestInsertSubjectSucceeds(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	id := getUUIDString()
	subjectTypeStr := getUUIDString()

	s := Subject{ID: id, Name: clusterName, Type: subjectTypeStr}
	gormDB := getGormHelper()
	gormDB.Create(&s)

	subject := Subject{}
	gormDB.First(&subject, "id = ?", id)

	expectedSubject := Subject{ID: id, Name: clusterName, Type: subjectTypeStr}
	assert.Equal(t, expectedSubject.ID, subject.ID, "expected %s got %s", expectedSubject.ID, subject.ID)
	assert.Equal(t, expectedSubject.Name, subject.Name, "expected %s got %s", expectedSubject.Name, subject.Name)
	assert.Equal(t, expectedSubject.Type, subject.Type, "expected %s got %s", expectedSubject.Type, subject.Type)

	// Drop the database instead of downgrading since we don't need the
	// data anyway
	if err := m.Drop(); err != nil {
		t.Fatalf("Unable to drop database: %s", err)
	}
}

func TestInsertSubjectWithLongNameFails(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	id := getUUIDString()
	maxNameLength := 256
	name := strings.Repeat("a", maxNameLength)
	subjectTypeStr := getUUIDString()

	s := Subject{ID: id, Name: name, Type: subjectTypeStr}
	gormDB := getGormHelper()
	err := gormDB.Create(&s).Error
	assert.NotEmpty(t, err, "Shouldn't be able to insert name values longer than 255 characters")
	// Drop the database instead of downgrading since we don't need the
	// data anyway
	if err := m.Drop(); err != nil {
		t.Fatalf("Unable to drop database: %s", err)
	}
}

func TestInsertSubjectWithNonUUIDFails(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	id := "1"
	subjectTypeStr := getUUIDString()

	s := Subject{ID: id, Name: clusterName, Type: subjectTypeStr}
	gormDB := getGormHelper()
	err := gormDB.Create(&s).Error
	fmt.Print(err)
	assert.NotEmpty(t, err, "Expect an error when creating IDs of the wrong type.")
	// Drop the database instead of downgrading since we don't need the
	// data anyway
	if err := m.Drop(); err != nil {
		t.Fatalf("Unable to drop database: %s", err)
	}
}

func TestInsertSubjectWithLongTypeFails(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	id := getUUIDString()
	maxTypeLength := 51
	subjectTypeStr := strings.Repeat("a", maxTypeLength)

	s := Subject{ID: id, Name: clusterName, Type: subjectTypeStr}
	gormDB := getGormHelper()
	err := gormDB.Create(&s).Error
	assert.NotEmpty(t, err, "Shouldn't be able to insert type values longer than 50 characters")
	// Drop the database instead of downgrading since we don't need the
	// data anyway
	if err := m.Drop(); err != nil {
		t.Fatalf("Unable to drop database: %s", err)
	}
}

func TestMigration(t *testing.T) { // nolint:paralleltest // database tests should run serially
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
