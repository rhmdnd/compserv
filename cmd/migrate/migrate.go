package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	config "github.com/rhmdnd/compserv/pkg/config"
)

func main() {
	var configDir = flag.String("config-dir", "configs/", "Path to YAML configuration directory containing a config.yaml file.")
	c := config.ParseConfig(*configDir)
	connStr := config.GetDatabaseConnectionString(c)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://migrations", c["db_name"], driver)
	if err != nil {
		log.Fatalf("Unable to load migrations: %s", err)
	}
	m.Up()
}
