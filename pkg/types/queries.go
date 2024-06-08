package types

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ClauseRepository struct {
	db *sql.DB
}

type ClauseController struct {
	repo *ClauseRepository
}

func NewClauseRepository(db *sql.DB) *ClauseRepository {
	return &ClauseRepository{db: db}
}

func NewClauseController(repo *ClauseRepository) *ClauseController {
	return &ClauseController{
		repo: repo,
	}
}

func (c *ClauseController) GetAllClauses(ctx *gin.Context) {
	clauses, err := c.repo.GetAllClauses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clauses"})
		return
	}
	ctx.JSON(http.StatusOK, clauses)
}

func (r *ClauseRepository) GetAllClauses() ([]Clause, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	query := "SELECT * FROM clause_section;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clauses := []Clause{}
	for rows.Next() {
		clause := Clause{}
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
