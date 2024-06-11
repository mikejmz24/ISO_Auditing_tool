package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"database/sql"
	"errors"
)

type Repository interface {
	GetAllClauses() ([]types.Clause, error)
	GetClauseByID(id int) (types.Clause, error)
	CreateClause(clause types.Clause) (int64, error)
	UpdateClause(clause types.Clause) error
	DeleteClause(id int) error
}

type repository struct {
	db *sql.DB
}

func NewClauseRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// GetAllClauses retrieves all clauses from the database
func (r *repository) GetAllClauses() ([]types.Clause, error) {
	query := "SELECT id, name, section FROM clause_section;"
	return executeQuery(r.db, query, scanClause)
}

// GetClauseByID retrieves a single clause by its ID
func (r *repository) GetClauseByID(id int) (types.Clause, error) {
	query := "SELECT id, name, section FROM clause_section WHERE id = ?;"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return types.Clause{}, err
	}
	defer rows.Close()

	if rows.Next() {
		clause, err := scanClause(rows)
		if err != nil {
			return types.Clause{}, err
		}
		return clause, nil
	}
	return types.Clause{}, errors.New("clause not found")
}

// CreateClause inserts a new clause into the database
func (r *repository) CreateClause(clause types.Clause) (int64, error) {
	query := "INSERT INTO clause_section (name, section) VALUES (?, ?);"
	result, err := r.db.Exec(query, clause.Name, clause.Section)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateClause updates an existing clause in the database
func (r *repository) UpdateClause(clause types.Clause) error {
	query := "UPDATE clause_section SET name = ?, section = ? WHERE id = ?;"
	_, err := r.db.Exec(query, clause.Name, clause.Section, clause.ID)
	return err
}

// DeleteClause deletes a clause from the database
func (r *repository) DeleteClause(id int) error {
	query := "DELETE FROM clause_section WHERE id = ?;"
	_, err := r.db.Exec(query, id)
	return err
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
