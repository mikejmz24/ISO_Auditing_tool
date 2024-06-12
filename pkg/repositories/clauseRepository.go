package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"database/sql"
	"encoding/json"
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
	query := `
		SELECT c.id, c.name, 
			IFNULL(JSON_ARRAYAGG(JSON_OBJECT('id', s.id, 'name', s.name)), '[]') AS sections 
		FROM clause c
		LEFT JOIN section s ON c.id = s.clause_id
		GROUP BY c.id, c.name;
	`
	return executeQuery(r.db, query, scanClause)
}

// GetClauseByID retrieves a single clause by its ID
func (r *repository) GetClauseByID(id int) (types.Clause, error) {
	query := `
		SELECT c.id, c.name, 
			IFNULL(JSON_ARRAYAGG(JSON_OBJECT('id', s.id, 'name', s.name)), '[]') AS sections 
		FROM clause c
		LEFT JOIN section s ON c.id = s.clause_id
		WHERE c.id = ?
		GROUP BY c.id, c.name;
	`
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
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO clause (name) VALUES (?);"
	result, err := tx.Exec(query, clause.Name)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	clauseID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, section := range clause.Sections {
		query := "INSERT INTO section (clause_id, name) VALUES (?, ?);"
		_, err := tx.Exec(query, clauseID, section.Name)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return clauseID, nil
}

// UpdateClause updates an existing clause in the database
func (r *repository) UpdateClause(clause types.Clause) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE clause SET name = ? WHERE id = ?;"
	_, err = tx.Exec(query, clause.Name, clause.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = "DELETE FROM section WHERE clause_id = ?;"
	_, err = tx.Exec(query, clause.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, section := range clause.Sections {
		query := "INSERT INTO section (clause_id, name) VALUES (?, ?);"
		_, err := tx.Exec(query, clause.ID, section.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return err
}

// DeleteClause deletes a clause from the database
func (r *repository) DeleteClause(id int) error {
	query := "DELETE FROM clause WHERE id = ?;"
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
	var sectionsJSON string
	err := rows.Scan(&clause.ID, &clause.Name, &sectionsJSON)
	if err != nil {
		return clause, err
	}

	err = json.Unmarshal([]byte(sectionsJSON), &clause.Sections)
	return clause, err
}
