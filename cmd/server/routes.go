package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ISO_Auditing_Tool/cmd/internal/database"
	"ISO_Auditing_Tool/templates"
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
		templ.Handler(templates.HelloForm()).ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/landing", func(c *gin.Context) {
		templ.Handler(templates.Base()).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/hello", func(c *gin.Context) {
		templates.HelloWebHandler(c.Writer, c.Request)
	})

	api := r.Group("/api")
	{
		api.GET("/clauses", s.apiClauseController.GetAllClauses)
	}

	r.GET("/clauses", s.htmlClauseController.GetAllClauses)
	r.GET("/clauses/add", func(c *gin.Context) {
		templ.Handler(templates.AddClause()).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/clauses/add", s.htmlClauseController.AddClause)
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
