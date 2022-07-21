package tests

import "time"

type Subject struct {
	ID   string
	Name string
	Type string
}

type Metadata struct {
	ID          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     string
	Description string
}
