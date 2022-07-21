package tests // nolint:testpackage

import (
	"database/sql"
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
	parentIDValue, err := subject.ParentID.Value()
	assert.False(t, subject.ParentID.Valid)
	assert.Nil(t, err)
	assert.Nil(t, parentIDValue)
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
}

func TestMetadataMigration(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()
	tableName := "metadata"

	version := uint(1)
	if err := m.Migrate(version); err != nil {
		t.Fatalf("Unable to upgrade database to version %d: %s", version, err)
	}

	// Ensure the metadata table doesn't exist before the upgrade
	result := gormDB.Migrator().HasTable(tableName)
	assert.False(t, result, "Table exists prior to migration: %s", tableName)

	// Ensure the table exists after running the migration
	version = uint(2)
	if err := m.Migrate(version); err != nil {
		t.Fatalf("Unable to upgrade database to version %d: %s", version, err)
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
	version = uint(1)
	if err := m.Migrate(version); err != nil {
		t.Fatalf("Unable to upgrade database to version %d: %s", version, err)
	}
	result = gormDB.Migrator().HasTable(tableName)
	assert.False(t, result, "Table exists after downgrade: %s", tableName)
}

func TestInsertMetadataSucceeds(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()

	if err := m.Migrate(2); err != nil {
		t.Fatalf("Unable to upgrade database to version %d: %s", 2, err)
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

	e := Metadata{ID: id, CreatedAt: createdAt, UpdatedAt: updatedAt, Version: version, Description: description}
	assert.Equal(t, e.ID, a.ID, "expected %s got %s", e.ID, a.ID)
	assert.Equal(t, e.CreatedAt, a.CreatedAt, "expected %s got %s", e.CreatedAt, a.CreatedAt)
	assert.Equal(t, e.UpdatedAt, a.UpdatedAt, "expected %s got %s", e.UpdatedAt, a.UpdatedAt)
	assert.Equal(t, e.Version, a.Version, "expected %s got %s", e.Version, a.Version)
	assert.Equal(t, e.Description, a.Description, "expected %s got %s", e.Description, a.Description)
}

func TestInsertSubjectWithParent(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	// Create a subject to act as the parent
	parentID := getUUIDString()
	subjectTypeStr := getUUIDString()
	p := Subject{ID: parentID, Name: clusterName, Type: subjectTypeStr}
	gormDB := getGormHelper()
	gormDB.Create(&p)

	// Create another subject referencing the parent
	id := getUUIDString()
	subjectTypeStr = getUUIDString()
	s := Subject{ID: id, Name: clusterName, Type: subjectTypeStr, ParentID: sql.NullString{String: parentID, Valid: true}}
	gormDB.Create(&s)

	a := Subject{}
	gormDB.First(&a, "id = ?", id)

	e := Subject{ID: id, Name: clusterName, Type: subjectTypeStr, ParentID: sql.NullString{String: parentID, Valid: true}}
	assert.Equal(t, e.ID, a.ID, "expected %s got %s", e.ID, a.ID)
	assert.Equal(t, e.Name, a.Name, "expected %s got %s", e.Name, a.Name)
	assert.Equal(t, e.Type, a.Type, "expected %s got %s", e.Type, a.Type)
	assert.True(t, a.ParentID.Valid)
	assert.Equal(t, e.ParentID.String, a.ParentID.String, "expected %s got %s", e.ParentID.String, a.ParentID.String)
}

func TestDeleteSubjectWithParent(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	// Create a subject to act as the parent
	parentID := getUUIDString()
	subjectTypeStr := getUUIDString()
	p := Subject{ID: parentID, Name: clusterName, Type: subjectTypeStr}
	gormDB := getGormHelper()
	err := gormDB.Create(&p).Error
	assert.Nil(t, err)

	// Create another subject referencing the parent
	id := getUUIDString()
	subjectTypeStr = getUUIDString()
	s := Subject{ID: id, Name: clusterName, Type: subjectTypeStr, ParentID: sql.NullString{String: parentID, Valid: true}}
	err = gormDB.Create(&s).Error
	assert.Nil(t, err)

	a := Subject{}
	gormDB.First(&a, "id = ?", id)

	e := Subject{ID: id, Name: clusterName, Type: subjectTypeStr, ParentID: sql.NullString{String: parentID, Valid: true}}
	assert.Equal(t, e.ID, a.ID, "expected %s got %s", e.ID, a.ID)
	assert.Equal(t, e.Name, a.Name, "expected %s got %s", e.Name, a.Name)
	assert.Equal(t, e.Type, a.Type, "expected %s got %s", e.Type, a.Type)
	assert.True(t, a.ParentID.Valid)
	assert.Equal(t, e.ParentID.String, a.ParentID.String, "expected %s got %s", e.ParentID.String, a.ParentID.String)

	err = gormDB.Delete(&p).Error
	assert.NotEmpty(t, err, "Deleting the parent subject should fail if the child subject still exists")

	// There should still be two subjects in the table
	var subjects []Subject
	result := gormDB.Find(&subjects)
	expectedSubjects := int64(2)
	assert.Equal(t, expectedSubjects, result.RowsAffected, "expected %d got %d", expectedSubjects, result.RowsAffected)

	// Deleting the child subject first should allow us to delete the
	// parent subject after since the foreign key isn't violated
	err = gormDB.Delete(&s).Error
	assert.Nil(t, err, "Failed to delete the child subject")
	err = gormDB.Delete(&p).Error
	assert.Nil(t, err, "Failed to delete the parent subject")

	result = gormDB.Find(&subjects)
	expectedSubjects = int64(0)
	assert.Equal(t, expectedSubjects, result.RowsAffected, "expected %d got %d", expectedSubjects, result.RowsAffected)
}

func TestInsertSubjectWithNonExistentParent(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	// Generate a fake parent ID that doesn't actually exist
	parentID := getUUIDString()

	// Create another subject referencing the parent
	id := getUUIDString()
	subjectTypeStr := getUUIDString()
	s := Subject{ID: id, Name: clusterName, Type: subjectTypeStr, ParentID: sql.NullString{String: parentID, Valid: true}}
	gormDB := getGormHelper()
	err := gormDB.Create(&s).Error
	assert.NotEmpty(t, err, "Shouldn't be able to insert values that violate foreign key constraints")
}

func TestInsertSubjectWithMetadata(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	if err := m.Up(); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	metadataID, err := insertMetadata()
	if err != nil {
		t.Fatalf("Unable to create metadata: %s", err)
	}

	// Create another subject referencing the parent
	id := getUUIDString()
	subjectTypeStr := getUUIDString()
	s := Subject{
		ID: id, Name: clusterName, Type: subjectTypeStr,
		MetadataID: sql.NullString{String: metadataID, Valid: true},
	}
	gormDB := getGormHelper()
	err = gormDB.Create(&s).Error
	assert.Nil(t, err, "Unable to create subject with metadata")
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
	expectedVersion = uint(9)
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

func TestSubjectParentIDMigration(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()

	type subject struct{}
	tableName := "subjects"
	columnName := "parent_id"

	version := uint(2)
	if err := m.Migrate(version); err != nil {
		t.Fatalf("Unable to migrate to version %d: %s", version, err)
	}
	result := gormDB.Migrator().HasTable(tableName)
	assert.True(t, result, "Table doesn't exist: %s", tableName)

	for _, s := range []string{"id", "name", "type"} {
		result = gormDB.Migrator().HasColumn(&subject{}, s)
		assert.True(t, result, "Table doesn't have column: %s", s)
	}
	result = gormDB.Migrator().HasColumn(&subject{}, columnName)
	assert.False(t, result, "Table has column: %s", columnName)

	// Ensure the upgrade adds the parent_id column and the constraint
	version = uint(3)
	if err := m.Migrate(version); err != nil {
		t.Fatalf("Unable to migrate to version %d: %s", version, err)
	}
	result = gormDB.Migrator().HasTable(tableName)
	assert.True(t, result, "Table doesn't exist: %s", tableName)

	for _, s := range []string{"id", "name", "type", "parent_id"} {
		result = gormDB.Migrator().HasColumn(&subject{}, s)
		assert.True(t, result, "Table doesn't have column: %s", s)
	}

	constraintName := "fk_subjects_parent_id"
	result = gormDB.Migrator().HasConstraint(&subject{}, constraintName)
	assert.True(t, result, "Table doesn't have constraint: %s", constraintName)

	// Make sure the downgrade removes the parent_id column and the
	// constraint
	version = uint(2)
	if err := m.Migrate(version); err != nil {
		t.Fatalf("Unable to migrate to version %d: %s", 2, err)
	}
	result = gormDB.Migrator().HasColumn(&subject{}, columnName)
	assert.False(t, result, "Table has column: %s", columnName)

	result = gormDB.Migrator().HasConstraint(&subject{}, constraintName)
	assert.False(t, result, "Table has constraint: %s", constraintName)
}

func TestSubjectMetadataIDMigration(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()

	type subject struct{}
	tableName := "subjects"
	columnName := "metadata_id"

	version := uint(4)
	if err := m.Migrate(version); err != nil {
		t.Fatalf("Unable to migrate to version %d: %s", version, err)
	}
	result := gormDB.Migrator().HasTable(tableName)
	assert.True(t, result, "Table doesn't exist: %s", tableName)

	for _, s := range []string{"id", "name", "type", "parent_id"} {
		result = gormDB.Migrator().HasColumn(&subject{}, s)
		assert.True(t, result, "Table doesn't have column: %s", s)
	}
	result = gormDB.Migrator().HasColumn(&subject{}, columnName)
	assert.False(t, result, "Table has column: %s", columnName)

	// Ensure the upgrade adds the metadata_id column and the constraint
	version = uint(5)
	if err := m.Migrate(version); err != nil {
		t.Fatalf("Unable to migrate to version %d: %s", version, err)
	}
	result = gormDB.Migrator().HasTable(tableName)
	assert.True(t, result, "Table doesn't have column: %s", tableName)

	for _, s := range []string{"id", "name", "type", "parent_id", "metadata_id"} {
		result = gormDB.Migrator().HasColumn(&subject{}, s)
		assert.True(t, result, "Table doesn't have column: %s", s)
	}

	constraintName := "fk_subjects_metadata_id"
	result = gormDB.Migrator().HasConstraint(&subject{}, constraintName)
	assert.True(t, result, "Table doesn't have constraint: %s", constraintName)

	// Make sure the downgrade removes the metadata_id column and the
	// constraint
	version = uint(4)
	if err := m.Migrate(version); err != nil {
		t.Fatalf("Unable to migrate to version %d: %s", version, err)
	}
	result = gormDB.Migrator().HasColumn(&subject{}, columnName)
	assert.False(t, result, "Table has column: %s", columnName)

	result = gormDB.Migrator().HasConstraint(&subject{}, constraintName)
	assert.False(t, result, "Table has constraint: %s", constraintName)
}

func TestAssessmentMigration(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()

	type assessments struct{}
	tableName := "assessments"

	if err := m.Migrate(3); err != nil {
		t.Fatalf("Unable to migrate to version %d: %s", 1, err)
	}
	result := gormDB.Migrator().HasTable(tableName)
	assert.False(t, result, "Table exist: %s", tableName)

	if err := m.Migrate(4); err != nil {
		t.Fatalf("Unable to migrate to version %d: %s", 4, err)
	}
	result = gormDB.Migrator().HasTable(tableName)
	assert.True(t, result, "Table doesn't exist: %s", tableName)

	for _, s := range []string{"id", "name", "metadata_id"} {
		result = gormDB.Migrator().HasColumn(&assessments{}, s)
		assert.True(t, result, "Table doesn't have column: %s", s)
	}

	constraintName := "fk_assessments_metadata_id"
	result = gormDB.Migrator().HasConstraint(&assessments{}, constraintName)
	assert.True(t, result, "Table doesn't have constraint: %s", constraintName)

	if err := m.Migrate(3); err != nil {
		t.Fatalf("Unable to migrate to version %d: %s", 3, err)
	}
	result = gormDB.Migrator().HasTable(tableName)
	assert.False(t, result, "Table exist: %s", tableName)
}

func TestResultMigration(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()
	tableName := "results"

	if err := m.Migrate(8); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	// Ensure the metadata table doesn't exist before the upgrade
	result := gormDB.Migrator().HasTable(tableName)
	assert.False(t, result, "Table exists prior to migration: %s", tableName)

	// Ensure the table exists after running the migration
	if err := m.Migrate(9); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	result = gormDB.Migrator().HasTable(tableName)
	assert.True(t, result, "Table doesn't exist: %s", tableName)

	// Check to make sure it has the expected columns
	type results struct{}
	columns := []string{
		"id", "name", "outcome", "instruction",
		"rationale", "control_id", "metadata_id", "subject_id",
		"assessment_id",
	}
	for _, s := range columns {
		result = gormDB.Migrator().HasColumn(&results{}, s)
		assert.True(t, result, "Column doesn't exist: %s", s)
	}

	constraintName := "fk_results_control_id"
	result = gormDB.Migrator().HasConstraint(&results{}, constraintName)
	assert.True(t, result, "Table doesn't have constraint: %s", constraintName)

	constraintName = "fk_results_metadata_id"
	result = gormDB.Migrator().HasConstraint(&results{}, constraintName)
	assert.True(t, result, "Table doesn't have constraint: %s", constraintName)

	constraintName = "fk_results_subject_id"
	result = gormDB.Migrator().HasConstraint(&results{}, constraintName)
	assert.True(t, result, "Table doesn't have constraint: %s", constraintName)

	constraintName = "fk_results_assessment_id"
	result = gormDB.Migrator().HasConstraint(&results{}, constraintName)
	assert.True(t, result, "Table doesn't have constraint: %s", constraintName)

	// Ensure the table is removed on downgrade
	if err := m.Migrate(8); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}
	result = gormDB.Migrator().HasTable(tableName)
	assert.False(t, result, "Table exists after downgrade: %s", tableName)
}

func TestInsertAssessmentSucceeds(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()

	if err := m.Migrate(4); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}
	id := getUUIDString()
	metadataID, err := insertMetadata()
	if err != nil {
		t.Fatalf("Unable to create metadata: %s", err)
	}
	name := getUUIDString()

	assessment := Assessment{ID: id, Name: name, MetadataID: metadataID}
	err = gormDB.Create(&assessment).Error
	assert.Nil(t, err)

	a := Assessment{}
	gormDB.First(&a, "id = ?", id)

	e := Assessment{ID: id, Name: name, MetadataID: metadataID}
	assert.Equal(t, e.ID, a.ID, "expected %s got %s", e.ID, a.ID)
	assert.Equal(t, e.Name, a.Name, "expected %s got %s", e.Name, a.Name)
	assert.Equal(t, e.MetadataID, a.MetadataID, "expected %s got %s", e.MetadataID, a.MetadataID)
}

func TestInsertResultSucceeds(t *testing.T) { // nolint:paralleltest // database tests should run serially
	m := getMigrationHelper(t)
	gormDB := getGormHelper()

	if err := m.Migrate(9); err != nil {
		t.Fatalf("Unable to upgrade database: %s", err)
	}

	metadataID, err := insertMetadata()
	if err != nil {
		t.Fatalf("Unable to create necessary metadata: %s", err)
	}
	subjectID, err := insertSubject()
	if err != nil {
		t.Fatalf("Unable to create necessary subject: %s", err)
	}
	controlID, err := insertControl()
	if err != nil {
		t.Fatalf("Unable to create necessary control: %s", err)
	}
	assessmentID, err := insertAssessment()
	if err != nil {
		t.Fatalf("Unable to create necessary assessment: %s", err)
	}

	id := getUUIDString()
	name := getUUIDString()
	outcome := getUUIDString()
	instruction := strings.Repeat(getUUIDString(), 100)
	rationale := getUUIDString()

	r := Result{
		ID: id, Name: name, ControlID: controlID,
		Outcome: outcome, Instruction: instruction,
		Rationale: rationale, MetadataID: metadataID,
		SubjectID: subjectID, AssessmentID: assessmentID,
	}

	err = gormDB.Create(&r).Error
	assert.Nil(t, err)

	a := Result{}
	gormDB.First(&a, "id = ?", id)

	e := Result{
		ID: id, Name: name, ControlID: controlID,
		Outcome: outcome, Instruction: instruction,
		Rationale: rationale, MetadataID: metadataID,
		SubjectID: subjectID, AssessmentID: assessmentID,
	}

	assert.Equal(t, e.ID, a.ID, "expected %s got %s", e.ID, a.ID)
	assert.Equal(t, e.Name, a.Name, "expected %s got %s", e.Name, a.Name)
	assert.Equal(t, e.ControlID, a.ControlID, "expected %s got %s", e.ControlID, a.ControlID)
	assert.Equal(t, e.Outcome, a.Outcome, "expected %s got %s", e.Outcome, a.Outcome)
	assert.Equal(t, e.Instruction, a.Instruction, "expected %s got %s", e.Instruction, a.Instruction)
	assert.Equal(t, e.Rationale, a.Rationale, "expected %s got %s", e.Rationale, a.Rationale)
	assert.Equal(t, e.MetadataID, a.MetadataID, "expected %s got %s", e.MetadataID, a.MetadataID)
	assert.Equal(t, e.SubjectID, a.SubjectID, "expected %s got %s", e.SubjectID, a.SubjectID)
	assert.Equal(t, e.AssessmentID, a.AssessmentID, "expected %s got %s", e.AssessmentID, a.AssessmentID)
}
