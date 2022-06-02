package e2e

import (
	"database/sql"
	"flag"
	"fmt"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var dbHost = flag.String("host", "", "Database host")
var dbPort = flag.String("port", "5432", "Database port")
var dbUsername = flag.String("username", "postgres", "Database username")
var dbPassword = flag.String("password", "", "Database password")
var dbName = flag.String("database", "compliance", "Database name")

func setUpDatabase(t *testing.T) *sql.DB {
	flag.Parse()
	// FIXME(rhmdnd): Find a better way to deal with the sslmode.
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", *dbHost, *dbUsername, *dbName, *dbPassword, *dbPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skip("Unable to initialize connection to test database")
	}
	if err := db.Ping(); err != nil {
		t.Skip("Unable to establish connection to test database")
	}

	// FIXME(rhmdnd): Finish adding migration logic.
	return db
}

func TestSubject(t *testing.T) {
	setUpDatabase(t)
}
