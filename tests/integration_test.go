package tests // nolint:testpackage

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Each test assumes the database is unmanaged. The test is responsible for
// setting up the state it requires for its test logic. This keeps the
// getMigrationHelper() method clean of any assumptions about what the tests
// expect from it. Each test must be run serially they have the power to change
// database schema and affect other tests.
func TestInsertSubjectSucceeds(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	id := getUUIDString()
	subjectTypeStr := getUUIDString()

	s := Subject{ID: id, Name: clusterName, Type: subjectTypeStr}
	gormDB := getGormHelper()
	gormDB.Create(&s)

	subject := Subject{}
	gormDB.First(&subject, "id = ?", id)

	expectedSubject := Subject{ID: id, Name: clusterName, Type: subjectTypeStr}
	assert.Equal(t, expectedSubject.ID, subject.ID, "expected %s got %s", expectedSubject.ID, subject.ID)
	assert.Equal(t, expectedSubject.Name, subject.Name, "expected %s got %s", expectedSubject.Name, subject.Name)
	assert.Equal(t, expectedSubject.Type, subject.Type, "expected %s got %s", expectedSubject.Type, subject.Type)

	// Drop the database instead of downgrading since we don't need the
	// data anyway
	if err := m.Drop(); err != nil {
		t.Fatalf("Unable to drop database: %s", err)
	}
}

func TestInsertSubjectWithLongNameFails(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	id := getUUIDString()
	maxNameLength := 256
	name := strings.Repeat("a", maxNameLength)
	subjectTypeStr := getUUIDString()

	s := Subject{ID: id, Name: name, Type: subjectTypeStr}
	gormDB := getGormHelper()
	err := gormDB.Create(&s).Error
	assert.NotEmpty(t, err, "Shouldn't be able to insert name values longer than 255 characters")
	// Drop the database instead of downgrading since we don't need the
	// data anyway
	if err := m.Drop(); err != nil {
		t.Fatalf("Unable to drop database: %s", err)
	}
}

func TestInsertSubjectWithNonUUIDFails(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	id := "1"
	subjectTypeStr := getUUIDString()

	s := Subject{ID: id, Name: clusterName, Type: subjectTypeStr}
	gormDB := getGormHelper()
	err := gormDB.Create(&s).Error
	fmt.Print(err)
	assert.NotEmpty(t, err, "Expect an error when creating IDs of the wrong type.")
	// Drop the database instead of downgrading since we don't need the
	// data anyway
	if err := m.Drop(); err != nil {
		t.Fatalf("Unable to drop database: %s", err)
	}
}

func TestInsertSubjectWithLongTypeFails(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	id := getUUIDString()
	maxTypeLength := 51
	subjectTypeStr := strings.Repeat("a", maxTypeLength)

	s := Subject{ID: id, Name: clusterName, Type: subjectTypeStr}
	gormDB := getGormHelper()
	err := gormDB.Create(&s).Error
	assert.NotEmpty(t, err, "Shouldn't be able to insert type values longer than 50 characters")
	// Drop the database instead of downgrading since we don't need the
	// data anyway
	if err := m.Drop(); err != nil {
		t.Fatalf("Unable to drop database: %s", err)
	}
}

func TestMigration(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)

	version, dirty, err := m.Version()

	expectedVersion := uint(0)
	expectedState := false
	assert.Equal(t, expectedVersion, version, "Database version mismatch: want %d but got %d", expectedVersion, version)
	assert.Equal(t, expectedState, dirty, "Database state mismatch: want %t but got %t", expectedState, dirty)
	// Currently, Version() doesn't return a typed error, but a generic one
	// with a specific string (e.g., "no migration"). If this changes, we
	// should update the test to check the error type returned and that
	// it's what we expect.
	assert.NotEmpty(t, err, "Collecting the version should return an error on an empty database")

	// Upgrade the database and make sure all upgrades apply cleanly.
	err = m.Up()
	version, dirty, _ = m.Version()
	expectedVersion = uint(1)
	assert.Equal(t, expectedVersion, version, "Database version mismatch: want %d but got %d", expectedVersion, version)
	assert.Equal(t, false, dirty, "Database state mismatch: want %t but got %t", false, dirty)
	assert.Equal(t, err, nil, "Error upgrading the database: %s", err)

	// Downgrade the database back to 0 and make sure all downgrades apply cleanly.
	err = m.Down()
	version, dirty, _ = m.Version()
	expectedVersion = uint(0)
	assert.Equal(t, expectedVersion, version, "Database version mismatch: want %d but got %d", expectedVersion, version)
	assert.Equal(t, false, dirty, "Database state mismatch: want %t but got %t", false, dirty)
	assert.Equal(t, err, nil, "Error downgrading the database: %s", err)
}
