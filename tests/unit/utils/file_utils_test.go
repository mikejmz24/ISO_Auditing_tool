package utils_test

import (
	"ISO_Auditing_Tool/pkg/utils"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestFileUtils struct {
	suite.Suite
}

// Helper function to check file discovery for 'up' or 'down' files.
func (suite *TestFileUtils) checkFilesForMigration(file string, direction string, expectedFiles []string) {
	root, _ := utils.GetProjectRoot()
	migrationsPath := filepath.Join(root, "internal", "migrations", "sql")

	// Test FindFilesInDir with explicit path
	foundFiles, err := utils.FindFilesInDir(migrationsPath, file, direction)
	assert.NoError(suite.T(), err, "FindFilesInDir() returned an error")

	// Normalize paths and ensure the found files match the expected ones
	var expectedFilesRes []string
	for _, expectedFile := range expectedFiles {
		expectedFilesRes = append(expectedFilesRes, filepath.Join(root, "internal", "migrations", "sql", expectedFile))
	}
	assert.Equal(suite.T(), expectedFilesRes, foundFiles, "Mismatched files. Got: %v, Expected: %v", foundFiles, expectedFilesRes)
}

func (suite *TestFileUtils) TestNoFileWithUp_ReturnsAllUpFiles() {
	output := []string{"001_base_tables.up.sql", "002_base_tables.up.sql"}
	suite.checkFilesForMigration("", "up", output)
}

func (suite *TestFileUtils) TestNoFileWithDown_ReturnsDownUpFiles() {
	output := []string{"001_base_tables.down.sql", "002_base_tables.down.sql"}
	suite.checkFilesForMigration("", "down", output)
}

func (suite *TestFileUtils) TestWithFileWithUp_ReturnsProvidedUpFileOnly() {
	output := []string{"001_base_tables.up.sql"}
	suite.checkFilesForMigration("001_base_tables", "up", output)
}

func (suite *TestFileUtils) TestWithFileWithDown_ReturnsProvidedDownFileOnly() {
	output := []string{"001_base_tables.down.sql"}
	suite.checkFilesForMigration("001_base_tables", "down", output)
}

func TestFileUtilsMethods(t *testing.T) {
	suite.Run(t, new(TestFileUtils))
}
