package seeds

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type SeedConfig struct {
	TableName   string
	CSVPath     string
	Columns     []string
	ForeignKeys []string // New field for tracking dependencies
	BatchSize   int      // New field for configurable batch size
}

var seedConfigs = []SeedConfig{
	{
		TableName: "reference_types",
		CSVPath:   filepath.Join("internal", "seeds", "csv", "reference_types.csv"),
		Columns:   []string{"name", "description"},
		BatchSize: 30,
	},
	{
		TableName:   "reference_values",
		CSVPath:     filepath.Join("internal", "seeds", "csv", "reference_values.csv"),
		Columns:     []string{"type_id", "code", "name", "description"},
		ForeignKeys: []string{"type_id"},
		BatchSize:   30,
	},
	{
		TableName: "standards",
		CSVPath:   filepath.Join("internal", "seeds", "csv", "standards.csv"),
		Columns:   []string{"name", "description", "version"},
		BatchSize: 3,
	},
	{
		TableName:   "requirement_level",
		CSVPath:     filepath.Join("internal", "seeds", "csv", "requirement_level.csv"),
		Columns:     []string{"standard_id", "level_name", "level_order"},
		ForeignKeys: []string{"standard_id"},
		BatchSize:   30,
	},
	{
		TableName:   "requirement",
		CSVPath:     filepath.Join("internal", "seeds", "csv", "requirements.csv"),
		Columns:     []string{"standard_id", "level_id", "parent_id", "reference_code", "name", "description"},
		ForeignKeys: []string{"standard_id"},
		BatchSize:   500,
	},
	{
		TableName:   "questions",
		CSVPath:     filepath.Join("internal", "seeds", "csv", "question.csv"),
		Columns:     []string{"requirement_id", "question", "guidance"},
		ForeignKeys: []string{"requirement_id"},
		BatchSize:   500,
	},
	{
		TableName:   "evidence",
		CSVPath:     filepath.Join("internal", "seeds", "csv", "evidence.csv"),
		Columns:     []string{"question_id", "type_id", "expected"},
		ForeignKeys: []string{"question_id", "type_id"},
		BatchSize:   500,
	},
}

// buildInsertQuery generates a parameterized query for batch inserts
func buildInsertQuery(tableName string, columns []string, batchSize int) string {
	placeholders := make([]string, batchSize)

	// For each row in the batch, create a tuple of placeholders
	singleRowPlaceholders := fmt.Sprintf("(%s)", strings.Join(strings.Split(strings.Repeat("?", len(columns)), ""), ","))
	for i := 0; i < batchSize; i++ {
		placeholders[i] = singleRowPlaceholders
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(columns, ","),
		strings.Join(placeholders, ","),
	)
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

	// Process records in batches
	totalRecords := len(records)
	for i := 0; i < totalRecords; i += config.BatchSize {
		end := i + config.BatchSize
		if end > totalRecords {
			end = totalRecords
		}

		batchSize := end - i
		if err := insertBatch(db, config, records[i:end], batchSize); err != nil {
			return fmt.Errorf("failed to insert batch for %s: %w", config.TableName, err)
		}
	}

	log.Printf("%v table seeded successfully", config.TableName)
	return nil
}

func insertBatch(db *sql.DB, config SeedConfig, records [][]string, batchSize int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Transaction rollback failed: %v", rbErr)
			}
		}
	}()

	// Build a query for the current batch size
	query := buildInsertQuery(config.TableName, config.Columns, batchSize)

	// Flatten all values into a single slice for the prepared statement
	args := make([]interface{}, 0, batchSize*len(config.Columns))
	for _, record := range records {
		for _, v := range record {
			if v == "" {
				args = append(args, nil)
			} else {
				args = append(args, v)
			}
		}
	}

	// Execute the batch insert
	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute batch statement: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

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

// ValidateData checks for referential integrity in the CSV data
// func ValidateData(db *sql.DB) error {
// 	for _, config := range seedConfigs {
// 		if len(config.ForeignKeys) > 0 {
// 			if err := validateForeignKeys(db, config); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
//
// func validateForeignKeys(db *sql.DB, config SeedConfig) error {
// 	// Implementation would check that referenced IDs exist
// 	// in parent tables before insertion
// 	return nil
// }
//
// Post-Seeding Verification Functions

// VerifySeeding runs checks on seeded data to ensure it was loaded correctly
func VerifySeeding(db *sql.DB) error {
	for _, config := range seedConfigs {
		if err := verifyTableCount(db, config); err != nil {
			return err
		}
	}
	return nil
}

func verifyTableCount(db *sql.DB, config SeedConfig) error {
	records, err := readCSV(config.CSVPath)
	if err != nil {
		return err
	}

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", config.TableName)
	if err := db.QueryRow(query).Scan(&count); err != nil {
		return err
	}

	if count != len(records) {
		return fmt.Errorf("table %s has %d rows but expected %d",
			config.TableName, count, len(records))
	}

	log.Printf("Verified %s has correct row count: %d", config.TableName, count)
	return nil
}

// TableInfo defines a table's configuration for seeding and truncation
type TableInfo struct {
	Name      string
	Priority  int      // Higher number = truncated later, seeded earlier
	DependsOn []string // Tables this table depends on
}

// Database tables configuration - ordered by dependency
var tableInfo = []TableInfo{
	{Name: "reference_values", Priority: 10},
	{Name: "reference_types", Priority: 20},
	{Name: "requirement_level", Priority: 30},
	{Name: "requirement", Priority: 40},
	{Name: "evidence", Priority: 50},
	{Name: "questions", Priority: 60},
	{Name: "standards", Priority: 70},
}

// TruncateOptions provides configuration for truncation operations
type TruncateOptions struct {
	DisableForeignKeyChecks bool
	Tables                  []string // Specific tables to truncate, empty means all
	VerifyEmpty             bool     // Verify tables are empty after truncation
	Timeout                 time.Duration
}

// DefaultTruncateOptions returns the default configuration for truncation
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		DisableForeignKeyChecks: true,
		Tables:                  nil, // All tables
		VerifyEmpty:             true,
		Timeout:                 30 * time.Second,
	}
}

// Truncate removes all data from specified tables
func Truncate(db *sql.DB, options TruncateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	log.Println("Truncating database tables...")
	startTime := time.Now()

	// Begin a transaction for the entire operation
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Transaction rollback failed: %v", rbErr)
			}
		}
	}()

	// Disable foreign key checks if requested
	if options.DisableForeignKeyChecks {
		if _, err := tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0"); err != nil {
			return fmt.Errorf("failed to disable foreign key checks: %w", err)
		}
		defer func() {
			// Re-enable foreign key checks regardless of truncation outcome
			if _, reEnableErr := db.Exec("SET FOREIGN_KEY_CHECKS = 1"); reEnableErr != nil {
				log.Printf("WARNING: Failed to re-enable foreign key checks: %v", reEnableErr)
			}
		}()
	}

	// Filter tables based on options
	tablesToTruncate := filterTables(tableInfo, options.Tables)
	if len(tablesToTruncate) == 0 {
		return fmt.Errorf("no tables to truncate")
	}

	// Truncate tables in priority order (lower priority first)
	sortTablesByPriority(tablesToTruncate)
	for _, table := range tablesToTruncate {
		tableStartTime := time.Now()

		query := fmt.Sprintf("TRUNCATE TABLE %s", table.Name)
		log.Printf("Executing: %s", query)

		if _, err := tx.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table.Name, err)
		}

		log.Printf("Truncated %s (took %v)", table.Name, time.Since(tableStartTime))
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Verify tables are empty if requested
	if options.VerifyEmpty {
		if err := verifyTablesEmpty(db, tablesToTruncate); err != nil {
			return fmt.Errorf("verification failed: %w", err)
		}
	}

	log.Printf("All tables truncated successfully (total time: %v)", time.Since(startTime))
	return nil
}

// filterTables returns table configs for the specified tables or all if none specified
func filterTables(configs []TableInfo, targetTables []string) []TableInfo {
	if len(targetTables) == 0 {
		return configs
	}

	// Convert target tables to a lookup map
	targetMap := make(map[string]bool)
	for _, t := range targetTables {
		targetMap[t] = true
	}

	// Filter configs to only include specified tables
	filteredConfigs := make([]TableInfo, 0, len(targetTables))
	for _, config := range configs {
		if targetMap[config.Name] {
			filteredConfigs = append(filteredConfigs, config)
		}
	}

	return filteredConfigs
}

// sortTablesByPriority sorts tables by priority (ascending)
func sortTablesByPriority(tables []TableInfo) {
	sort.Slice(tables, func(i, j int) bool {
		return tables[i].Priority < tables[j].Priority
	})
}

// verifyTablesEmpty checks that all truncated tables have zero rows
func verifyTablesEmpty(db *sql.DB, tables []TableInfo) error {
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table.Name)

		if err := db.QueryRow(query).Scan(&count); err != nil {
			return fmt.Errorf("failed to verify table %s is empty: %w", table.Name, err)
		}

		if count > 0 {
			return fmt.Errorf("table %s still has %d rows after truncation", table.Name, count)
		}

		log.Printf("Verified %s is empty", table.Name)
	}

	return nil
}

// RefreshDatabase truncates all tables and then reseeds them
func RefreshDatabase(db *sql.DB) error {
	// First truncate all tables
	if err := Truncate(db, DefaultTruncateOptions()); err != nil {
		return fmt.Errorf("truncation failed: %w", err)
	}

	// Then reseed them
	if err := Seed(db); err != nil {
		return fmt.Errorf("reseeding failed: %w", err)
	}

	log.Println("Database truncate and refresh completed successfully")
	return nil
}
