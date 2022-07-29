package tests

import (
	"database/sql"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Necessary to invoke migrations locally in the repository
	"github.com/google/uuid"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Assessments struct {
	ID         string
	Name       string
	MetadataID string
}

type Metadata struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     string
	Description string
}

type Subject struct {
	ID       string
	Name     string
	Type     string
	ParentID sql.NullString
}

func getDatabaseConnection(t *testing.T) *sql.DB {
	t.Helper()
	// Generlize this so that it can be used to connect to any Postgres
	// database to run tests.
	connStr := "host=localhost user=dbadmin dbname=compliance password=secret port=5432 sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skipf("Unable to initialize connection to test database: %s", err)
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
				t.Skipf("Unable to establish connection to test database: %s", err)
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
		t.Skipf("Unable to initialize database driver for migrations: %s", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	if err != nil {
		t.Skipf("Unable to initialize migrations: %s", err)
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

func insertMetadata() (string, error) {
	gormDB := getGormHelper()

	id := getUUIDString()
	createdAt := time.Now().UTC().Round(time.Microsecond)
	updatedAt := time.Now().UTC().Round(time.Microsecond)
	version := getUUIDString()
	description := getUUIDString()

	md := Metadata{ID: id, CreatedAt: createdAt, UpdatedAt: updatedAt, Version: version, Description: description}
	if err := gormDB.Create(&md).Error; err != nil {
		return "", err
	}

	return id, nil
}
