package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// FindFilesInDir handles discovery of migration files based on parameters
func FindFilesInDir(file string, direction string) ([]string, error) {
	root, err := GetProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to get project root: %w", err)
	}

	migrationsPath := filepath.Join(root, "cmd", "internal", "migrations", "sql")
	var files []string

	// Handle specific file case
	if file != "" {
		filePath := filepath.Join(migrationsPath, fmt.Sprintf("%s.%s.sql", file, direction))
		if matches, _ := filepath.Glob(filePath); len(matches) == 0 {
			return nil, fmt.Errorf("migration file not found: %s", filePath)
		}
		files = []string{filePath}
	} else {
		// Get all migration files for the specified direction
		log.Printf("Direction received: %s", direction)
		pattern := filepath.Join(migrationsPath, fmt.Sprintf("*.%s.sql", direction))
		log.Printf("Searching for migrations with pattern: %s", pattern)

		files, err = filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to find migration files: %w", err)
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("no migration files found in %s", migrationsPath)
		}
		// Sort files to ensure consistent order
		sort.Strings(files)
	}

	return files, nil
}

func GetProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	for {
		// Check if go.mod exists in the current directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir { // Reached the root of the filesystem
			return "", fmt.Errorf("go.mod not found; ensure you're running within a Go project")
		}
		dir = parent
	}
}
