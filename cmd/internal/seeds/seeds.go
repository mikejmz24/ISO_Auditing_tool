package seeds

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type SeedConfig struct {
	TableName   string
	CSVPath     string
	InsertQuery string
}

var seedConfigs = []SeedConfig{
	{
		TableName:   "iso_standard",
		CSVPath:     filepath.Join("cmd", "internal", "seeds", "csv", "iso_standard.csv"),
		InsertQuery: "INSERT INTO iso_standard (name) VALUES (?)",
	},
	{
		TableName:   "clause",
		CSVPath:     filepath.Join("cmd", "internal", "seeds", "csv", "clause.csv"),
		InsertQuery: "INSERT INTO clause (id, iso_standard_id, name) VALUES (?, ?, ?)",
	},
	{
		TableName:   "section",
		CSVPath:     filepath.Join("cmd", "internal", "seeds", "csv", "section.csv"),
		InsertQuery: "INSERT INTO section (id, clause_id, name) VALUES (?, ?, ?)",
	},
	{
		TableName:   "subsection",
		CSVPath:     filepath.Join("cmd", "internal", "seeds", "csv", "subsection.csv"),
		InsertQuery: "INSERT INTO subsection (id, section_id, name) VALUES (?, ?, ?)",
	},
	{
		TableName:   "question",
		CSVPath:     filepath.Join("cmd", "internal", "seeds", "csv", "question.csv"),
		InsertQuery: "INSERT INTO question (id, section_id, subsection_id, name) VALUES (?, ?, ?, ?)",
	},
}

func Seed(db *sql.DB) error {
	for _, config := range seedConfigs {
		if err := seedTable(db, config); err != nil {
			log.Printf("Failed to execute seed function for table %v: %v", config.TableName, err)
			return fmt.Errorf("failed to execute seed: %w", err)
		}
	}
	log.Println("Database tables seeded successfully")
	return nil
}

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv data: %w", err)
	}
	return records, nil
}

func seedTable(db *sql.DB, config SeedConfig) error {
	log.Printf("Seeding %v table...", config.TableName)
	if isTableSeeded(db, config.TableName) {
		log.Printf("%v table already seeded, skipping...", config.TableName)
		return nil
	}

	records, err := readCSV(config.CSVPath)
	if err != nil {
		return fmt.Errorf("failed to read csv for %v: %w", config.TableName, err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.Prepare(config.InsertQuery)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Transaction rollback failed: %v", rollbackErr)
		}
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, record := range records {
		args := prepareArgs(record)
		if _, err := stmt.Exec(args...); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Transaction rollback failed: %v", rollbackErr)
			}
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Printf("Transaction rollback failed: %v", rollbackErr)
		}
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("%v table seeded successfully", config.TableName)
	return nil
}

func isTableSeeded(db *sql.DB, tableName string) bool {
	var count int
	queryString := fmt.Sprintf("SELECT COUNT(*) FROM %v", tableName)
	row := db.QueryRow(queryString)
	if err := row.Scan(&count); err != nil {
		log.Printf("Error counting rows in %v: %v", tableName, err)
		return false
	}
	return count > 0
}

func prepareArgs(record []string) []any {
	args := make([]any, len(record))
	for i, v := range record {
		if v == "" {
			args[i] = nil
		} else {
			args[i] = v
		}
	}
	return args
}
