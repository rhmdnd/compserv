package compserv

import (
	context "context"
	"log"

	"gorm.io/gorm"
)

type server struct {
	UnimplementedComplianceServiceServer
	database *gorm.DB
}

func NewServer(db *gorm.DB) *server { // nolint:revive,golint // returning a private struct from an exported fn is fine
	return &server{database: db}
}

func (s *server) SetResult(ctx context.Context, result *ResultRequest) (*ResultResponse, error) {
	log.Printf("%+v", result)
	// Wire this up to the database and persist the result
	return &ResultResponse{}, nil
}
