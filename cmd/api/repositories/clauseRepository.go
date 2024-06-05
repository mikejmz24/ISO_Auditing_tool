package repositories

import (
	"ISO_Auditing_Tool/cmd/api/types"
	"database/sql"
)

type ClauseRepository interface {
	GetAllClauses() ([]types.Clause, error)
}

type clauseRepository struct {
	db *sql.DB
}

func NewClauseRepository(db *sql.DB) ClauseRepository {
	return &clauseRepository{
		db: db,
	}
}

func (r *clauseRepository) GetAllClauses() ([]types.Clause, error) {
	query := "SELECT id, name, section FROM clause_section;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clauses := []types.Clause{}
	for rows.Next() {
		clause := types.Clause{}
		err := rows.Scan(&clause.ID, &clause.Name, &clause.Section)
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, clause)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clauses, nil
}
