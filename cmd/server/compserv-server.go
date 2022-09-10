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
	configDir := flag.String("config-dir", "configs/",
		"Path to YAML configuration directory containing a config.yaml file.")
	configFile := flag.String("config-file", "config.yaml",
		"File name of the service config")
	flag.Parse()
	v := config.ParseConfig(*configDir, *configFile)
	connStr := config.GetDatabaseConnectionString(v)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	log.Printf("Connected to database: %v", v.GetString("database.host"))

	appStr := v.GetString("app.host") + ":" + v.GetString("app.port")
	lis, err := net.Listen("tcp", appStr)
	if err != nil {
		log.Fatalf("Failed to listen to %s: %v", appStr, err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	api.RegisterComplianceServiceServer(grpcServer, api.NewServer(db))
	log.Printf("Server listening on %s", appStr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start grpc server %v", err)
	}
}
