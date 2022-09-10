package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"net"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Necessary to invoke migrations from files
	config "github.com/rhmdnd/compserv/pkg/config"
	"github.com/spf13/viper"
)

func main() {
	configDir := flag.String("config-dir", "configs/",
		"Path to YAML configuration directory containing a config.yaml file.")
	configFile := flag.String("config-file", "config.yaml",
		"File name of the service config")
	flag.Parse()
	v := config.ParseConfig(*configDir, *configFile)
	db := getDatabaseConnection(v)
	log.Printf("Connected to database: %v", v.GetString("database.host"))
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Unable to initialize database driver for migrations: %s", err)
	}
	// The file path to the migration could be a configuration option, but
	// I'm not sure how useful that would be since they're copied into the
	// container during the build process. Might only be useful for people
	// building their own container images.
	m, err := migrate.NewWithDatabaseInstance("file:///app/migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("Unable to initialize migrations: %s", err)
	}
	if err := m.Up(); err != nil {
		log.Fatalf("Unable to upgrade the database: %s", err)
	}
	version, _, err := m.Version()
	if err != nil {
		log.Fatalf("Unable to determine database version: %s", err)
	}
	log.Printf("Database successful migrated to version %d", version)
}

func getDatabaseConnection(v *viper.Viper) *sql.DB {
	// This should be updated so that we don't have to disable ssl
	connStr := config.GetDatabaseConnectionString(v) + " sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to initialize connection to database: %s", err)
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
				log.Fatalf("Retrying database connection due to error: %s", err)
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
				log.Fatalf("Unable to establish connection to database: %s", err)
			}
		}
	}

	return db
}
