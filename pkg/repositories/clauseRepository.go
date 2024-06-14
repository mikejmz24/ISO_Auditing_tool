package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"database/sql"
	"encoding/json"
	"errors"
)

// Repository interface defines the methods for interacting with the database
type Repository interface {
	// Clause methods
	GetAllClauses() ([]types.Clause, error)
	GetClauseByID(id int) (types.Clause, error)
	CreateClause(clause types.Clause) (int64, error)
	UpdateClause(clause types.Clause) error
	DeleteClause(id int) error

	// ISO Standard methods
	GetAllISOStandards() ([]types.ISOStandard, error)
	GetISOStandardByID(id int) (types.ISOStandard, error)
	CreateISOStandard(standard types.ISOStandard) (int64, error)
	UpdateISOStandard(standard types.ISOStandard) error
	DeleteISOStandard(id int) error

	// Section methods
	GetAllSections() ([]types.Section, error)
	GetSectionByID(id int) (types.Section, error)
	CreateSection(section types.Section) (int64, error)
	UpdateSection(section types.Section) error
	DeleteSection(id int) error

	// Question methods
	GetAllQuestions() ([]types.Question, error)
	GetQuestionByID(id int) (types.Question, error)
	CreateQuestion(question types.Question) (int64, error)
	UpdateQuestion(question types.Question) error
	DeleteQuestion(id int) error
}

// repository struct holds the database connection
type repository struct {
	db *sql.DB
}

// NewClauseRepository creates a new instance of the repository
func NewClauseRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// Clause methods

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

func (r *repository) CreateClause(clause types.Clause) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO clause (name, iso_standard_id) VALUES (?, ?);"
	result, err := tx.Exec(query, clause.Name, clause.ISOStandardID)
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

func (r *repository) UpdateClause(clause types.Clause) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := "UPDATE clause SET name = ?, iso_standard_id = ? WHERE id = ?;"
	_, err = tx.Exec(query, clause.Name, clause.ISOStandardID, clause.ID)
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

func (r *repository) DeleteClause(id int) error {
	query := "DELETE FROM clause WHERE id = ?;"
	_, err := r.db.Exec(query, id)
	return err
}

// ISO Standard methods

func (r *repository) GetAllISOStandards() ([]types.ISOStandard, error) {
	query := "SELECT id, name FROM iso_standard;"
	return executeQuery(r.db, query, scanISOStandard)
}

func (r *repository) GetISOStandardByID(id int) (types.ISOStandard, error) {
	query := "SELECT id, name FROM iso_standard WHERE id = ?;"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return types.ISOStandard{}, err
	}
	defer rows.Close()

	if rows.Next() {
		standard, err := scanISOStandard(rows)
		if err != nil {
			return types.ISOStandard{}, err
		}
		return standard, nil
	}
	return types.ISOStandard{}, errors.New("ISO standard not found")
}

func (r *repository) CreateISOStandard(standard types.ISOStandard) (int64, error) {
	query := "INSERT INTO iso_standard (name) VALUES (?);"
	result, err := r.db.Exec(query, standard.Name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *repository) UpdateISOStandard(standard types.ISOStandard) error {
	query := "UPDATE iso_standard SET name = ? WHERE id = ?;"
	_, err := r.db.Exec(query, standard.Name, standard.ID)
	return err
}

func (r *repository) DeleteISOStandard(id int) error {
	query := "DELETE FROM iso_standard WHERE id = ?;"
	_, err := r.db.Exec(query, id)
	return err
}

// Section methods

func (r *repository) GetAllSections() ([]types.Section, error) {
	query := "SELECT id, name, clause_id FROM section;"
	return executeQuery(r.db, query, scanSection)
}

func (r *repository) GetSectionByID(id int) (types.Section, error) {
	query := "SELECT id, name, clause_id FROM section WHERE id = ?;"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return types.Section{}, err
	}
	defer rows.Close()

	if rows.Next() {
		section, err := scanSection(rows)
		if err != nil {
			return types.Section{}, err
		}
		return section, nil
	}
	return types.Section{}, errors.New("section not found")
}

func (r *repository) CreateSection(section types.Section) (int64, error) {
	query := "INSERT INTO section (clause_id, name) VALUES (?, ?);"
	result, err := r.db.Exec(query, section.ClauseID, section.Name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *repository) UpdateSection(section types.Section) error {
	query := "UPDATE section SET name = ?, clause_id = ? WHERE id = ?;"
	_, err := r.db.Exec(query, section.Name, section.ClauseID, section.ID)
	return err
}

func (r *repository) DeleteSection(id int) error {
	query := "DELETE FROM section WHERE id = ?;"
	_, err := r.db.Exec(query, id)
	return err
}

// Question methods

func (r *repository) GetAllQuestions() ([]types.Question, error) {
	query := "SELECT id, text, section_id FROM question;"
	return executeQuery(r.db, query, scanQuestion)
}

func (r *repository) GetQuestionByID(id int) (types.Question, error) {
	query := "SELECT id, text, section_id FROM question WHERE id = ?;"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return types.Question{}, err
	}
	defer rows.Close()

	if rows.Next() {
		question, err := scanQuestion(rows)
		if err != nil {
			return types.Question{}, err
		}
		return question, nil
	}
	return types.Question{}, errors.New("question not found")
}

func (r *repository) CreateQuestion(question types.Question) (int64, error) {
	query := "INSERT INTO question (section_id, text) VALUES (?, ?);"
	result, err := r.db.Exec(query, question.SectionID, question.Text)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *repository) UpdateQuestion(question types.Question) error {
	query := "UPDATE question SET text = ?, section_id = ? WHERE id = ?;"
	_, err := r.db.Exec(query, question.Text, question.SectionID, question.ID)
	return err
}

func (r *repository) DeleteQuestion(id int) error {
	query := "DELETE FROM question WHERE id = ?;"
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

// scanISOStandard scans a single row into an ISOStandard struct
func scanISOStandard(rows *sql.Rows) (types.ISOStandard, error) {
	var standard types.ISOStandard
	err := rows.Scan(&standard.ID, &standard.Name)
	return standard, err
}

// scanSection scans a single row into a Section struct
func scanSection(rows *sql.Rows) (types.Section, error) {
	var section types.Section
	err := rows.Scan(&section.ID, &section.Name, &section.ClauseID)
	return section, err
}

// scanQuestion scans a single row into a Question struct
func scanQuestion(rows *sql.Rows) (types.Question, error) {
	var question types.Question
	err := rows.Scan(&question.ID, &question.Text, &question.SectionID)
	return question, err
}
