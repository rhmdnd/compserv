package main

import (
	"flag"
	"log"
	"net"

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

	log.Pritf("Connected to database: %v", db)

	appStr := c["app_host"] + ":" + c["app_port"]
	lis, err := net.Listen("tcp", appStr)
	if err != nil {
		log.Fatalf("Failed to listen to %s: %v", appStr, err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	api.RegisterComplianceServiceServer(grpcServer, api.NewServer(db))
	log.Printf("Server listening on %s", appStr)
	grpcServer.Serve(lis)
}
