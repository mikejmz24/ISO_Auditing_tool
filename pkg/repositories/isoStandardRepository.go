package repositories

import (
	"ISO_Auditing_Tool/pkg/types"
	"database/sql"
	"errors"
)

// Repository interface defines the methods for interacting with the database
type IsoStandardRepository interface {
	// ISO Standard methods
	GetAllISOStandards() ([]types.ISOStandard, error)
	GetISOStandardByID(id int) (types.ISOStandard, error)
	CreateISOStandard(standard types.ISOStandard) (int64, error)
	UpdateISOStandard(standard types.ISOStandard) error
	DeleteISOStandard(id int) error
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

func (r *isoStandardRepository) GetISOStandardByID(id int) (types.ISOStandard, error) {
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

func (r *isoStandardRepository) CreateISOStandard(standard types.ISOStandard) (int64, error) {
	query := "INSERT INTO iso_standard (name) VALUES (?);"
	result, err := r.db.Exec(query, standard.Name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *isoStandardRepository) UpdateISOStandard(standard types.ISOStandard) error {
	query := "UPDATE iso_standard SET name = ? WHERE id = ?;"
	_, err := r.db.Exec(query, standard.Name, standard.ID)
	return err
}

func (r *isoStandardRepository) DeleteISOStandard(id int) error {
	query := "DELETE FROM iso_standard WHERE id = ?;"
	_, err := r.db.Exec(query, id)
	return err
}
