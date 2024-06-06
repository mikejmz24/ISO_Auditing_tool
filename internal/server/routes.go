package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// "ISO_Auditing_Tool/cmd/api/types"
	"ISO_Auditing_Tool/cmd/web"
	"ISO_Auditing_Tool/internal/database"
	"database/sql"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes(db *sql.DB) http.Handler {
	r := gin.Default()
	s.db = database.New()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.Static("/assets", "./cmd/web/assets")

	r.GET("/web", func(c *gin.Context) {
		templ.Handler(web.HelloForm()).ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/landing", func(c *gin.Context) {
		templ.Handler(web.Base()).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/hello", func(c *gin.Context) {
		web.HelloWebHandler(c.Writer, c.Request)
	})

	api := r.Group("/api")
	{
		api.GET("/clauses", s.clauseController.GetAllClauses)
	}
	// api.GET("/clauses", s.clauseController.GetAllClauses)

	r.GET("/clauses", s.clauseController.GetAllClauses)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

// func (s *Server) getAllClausesHandler(c *gin.Context) {
// 	clauses, err := s.db.GetAllClauses()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": clauses})
// }
