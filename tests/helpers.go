package tests

import "database/sql"

type Subject struct {
	ID       string
	Name     string
	Type     string
	ParentID sql.NullString
}
