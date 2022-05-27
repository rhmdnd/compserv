package main

import (
	"flag"
	"log"

	api "github.com/rhmdnd/compserv/pkg/api"
	config "github.com/rhmdnd/compserv/pkg/config"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var configDir = flag.String("config-dir", "configs/", "Path to YAML configuration directory containing a config.yaml file.")
	flag.Parse()
	c := config.ParseConfig(*configDir)
	connStr := config.GetDatabaseConnectionString(c)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	log.Printf("Connected to database: %v", db)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	api.RegisterComplianceServiceServer(grpcServer, api.NewServer(db))
	// TODO(rhmdnd): Add host and port configuration then call
	// grpcServer.Server()
}
