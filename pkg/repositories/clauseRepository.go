package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"database/sql"
)

type Repository interface {
	GetAllClauses() ([]types.Clause, error)
}

type repository struct {
	db *sql.DB
}

func NewClauseRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAllClauses() ([]types.Clause, error) {
	query := "SELECT id, name, section FROM clause_section;"
	return executeQuery(r.db, query, scanClause)
}

// executeQuery is a generic function that executes a query and scans the results
func executeQuery[T any](db *sql.DB, query string, scanFunc func(*sql.Rows) (T, error)) ([]T, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		item, err := scanFunc(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// scanClause scans a single row into a Clause struct
func scanClause(rows *sql.Rows) (types.Clause, error) {
	var clause types.Clause
	err := rows.Scan(&clause.ID, &clause.Name, &clause.Section)
	return clause, err
}
