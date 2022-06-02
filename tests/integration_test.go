package tests // nolint:testpackage

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

func TestMetadataMigration(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()
	tableName := "metadata"

	if err := m.Migrate(1); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	// Ensure the metadata table doesn't exist before the upgrade
	result := gormDB.Migrator().HasTable(tableName)
	assert.False(t, result, "Table exists prior to migration: %s", tableName)

	// Ensure the table exists after running the migration
	if err := m.Migrate(2); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	result = gormDB.Migrator().HasTable(tableName)
	assert.True(t, result, "Table doesn't exist: %s", tableName)

	// Check to make sure it has the expected columns
	type metadata struct{}
	for _, s := range []string{"id", "created_at", "updated_at", "version", "description"} {
		result = gormDB.Migrator().HasColumn(&metadata{}, s)
		assert.True(t, result, "Column doesn't exist: %s", s)
	}

	// Ensure the table is removed on downgrade
	if err := m.Migrate(1); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}
	result = gormDB.Migrator().HasTable(tableName)
	assert.False(t, result, "Table exists after downgrade: %s", tableName)

	if err := m.Drop(); err != nil {
		t.Fatalf("Unable to drop database: %s", err)
	}
}

func TestInsertMetadataSucceeds(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()

	if err := m.Migrate(2); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	id := getUUIDString()
	// time.Now().UTC() will return a time.Time struct with microsecond
	// precision, but the timestamp type in PostgreSQL is not microsecond
	// precise. Round off the microsecond precision to model what the
	// database actually gives us.
	createdAt := time.Now().UTC().Round(time.Microsecond)
	updatedAt := time.Now().UTC().Round(time.Microsecond)
	version := getUUIDString()
	description := getUUIDString()

	md := Metadata{ID: id, CreatedAt: createdAt, UpdatedAt: updatedAt, Version: version, Description: description}
	err := gormDB.Create(&md).Error
	assert.Nil(t, err)

	a := Metadata{}
	gormDB.First(&a, "id = ?", id)
	fmt.Println(a.UpdatedAt)

	e := Metadata{ID: id, CreatedAt: createdAt, UpdatedAt: updatedAt, Version: version, Description: description}
	assert.Equal(t, e.ID, a.ID, "expected %s got %s", e.ID, a.ID)
	assert.Equal(t, e.CreatedAt, a.CreatedAt, "expected %s got %s", e.CreatedAt, a.CreatedAt)
	assert.Equal(t, e.UpdatedAt, a.UpdatedAt, "expected %s got %s", e.UpdatedAt, a.UpdatedAt)
	assert.Equal(t, e.Version, a.Version, "expected %s got %s", e.Version, a.Version)
	assert.Equal(t, e.Description, a.Description, "expected %s got %s", e.Description, a.Description)

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
	expectedVersion = uint(2)
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
