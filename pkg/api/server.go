package compserv

import "gorm.io/gorm"

type server struct {
	UnimplementedComplianceServiceServer
	database *gorm.DB
}

func NewServer(db *gorm.DB) *server { // nolint:revive,golint // returning a private struct from an exported fn is fine
	return &server{database: db}
}
