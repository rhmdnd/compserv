package compserv

import "gorm.io/gorm"

type server struct {
	UnimplementedComplianceServiceServer
	database *gorm.DB
}

func NewServer(db *gorm.DB) *server {
	return &server{database: db}
}
