package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"database/sql"
	"errors"
	"fmt"
)

// Repository interface defines the methods for interacting with the database
type IsoStandardRepository interface {
	// ISO Standard methods
	GetAllISOStandards() ([]types.ISOStandard, error)
	GetISOStandardByID(id int64) (types.ISOStandard, error)
	CreateISOStandard(isoStandard types.ISOStandard) (types.ISOStandard, error)
	UpdateISOStandard(isoStandard types.ISOStandard) error
	DeleteISOStandard(id int64) error
}

// repository struct holds the database connection
type isoStandardRepository struct {
	db *sql.DB
}

// NewClauseRepository creates a new instance of the repository
func NewIsoStandardRepository(db *sql.DB) IsoStandardRepository {
	return &isoStandardRepository{
		db: db,
	}
}

// ISO Standard methods
func (r *isoStandardRepository) GetAllISOStandards() ([]types.ISOStandard, error) {
	query := "SELECT id, name FROM iso_standard;"
	return executeQuery(r.db, query, scanISOStandard)
}

func (r *isoStandardRepository) GetISOStandardByID(id int64) (types.ISOStandard, error) {
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

func (r *isoStandardRepository) CreateISOStandard(standard types.ISOStandard) (types.ISOStandard, error) {
	query := "INSERT INTO iso_standard (name) VALUES (?);"
	result, err := r.db.Exec(query, standard.Name)
	if err != nil {
		return types.ISOStandard{}, err
	}
	// return result.LastInsertId()
	id, err := result.LastInsertId()
	if err != nil {
		return types.ISOStandard{}, err
	}
	standard.ID = int(id)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.ISOStandard{}, err
	}
	fmt.Printf("Created %v row(s)", rowsAffected)
	return standard, nil
}

func (r *isoStandardRepository) UpdateISOStandard(standard types.ISOStandard) error {
	query := "UPDATE iso_standard SET name = ? WHERE id = ?;"
	result, err := r.db.Exec(query, standard.Name, standard.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *isoStandardRepository) DeleteISOStandard(id int64) error {
	query := "DELETE FROM iso_standard WHERE id = ?;"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}
