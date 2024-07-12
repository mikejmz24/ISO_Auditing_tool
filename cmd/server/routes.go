package server

import (
	"net/http"

	"ISO_Auditing_Tool/pkg/middleware"
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

	// Web routes for rendering HTML
	r.GET("/web", func(c *gin.Context) {
		templ.Handler(templates.HelloForm()).ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/landing", func(c *gin.Context) {
		templ.Handler(templates.Base()).ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/hello", func(c *gin.Context) {
		templates.HelloWebHandler(c.Writer, c.Request)
	})

	// API routes group
	api := r.Group("/api")
	api.Use(middleware.ErrorHandler())
	{
		api.GET("/iso_standards", s.apiIsoStandardController.GetAllISOStandards)
		api.GET("/iso_standards/:id", s.apiIsoStandardController.GetISOStandardByID)
		api.POST("/iso_standards", s.apiIsoStandardController.CreateISOStandard)
		api.PUT("/iso_standards/:id", s.apiIsoStandardController.UpdateISOStandard)
		api.DELETE("/iso_standards/:id", s.apiIsoStandardController.DeleteISOStandard)

		api.GET("/clauses", s.apiClauseController.GetAllClauses)
		api.GET("/clauses/:id", s.apiClauseController.GetClauseByID)
		api.POST("/clauses", s.apiClauseController.CreateClause)
		api.PUT("/clauses/:id", s.apiClauseController.UpdateClause)
		api.DELETE("/clauses/:id", s.apiClauseController.DeleteClause)

		api.GET("/sections", s.apiClauseController.GetAllSections)
		api.GET("/sections/:id", s.apiClauseController.GetSectionByID)
		api.POST("/sections", s.apiClauseController.CreateSection)
		api.PUT("/sections/:id", s.apiClauseController.UpdateSection)
		api.DELETE("/sections/:id", s.apiClauseController.DeleteSection)

		api.GET("/questions", s.apiClauseController.GetAllQuestions)
		api.GET("/questions/:id", s.apiClauseController.GetQuestionByID)
		api.POST("/questions", s.apiClauseController.CreateQuestion)
		api.PUT("/questions/:id", s.apiClauseController.UpdateQuestion)
		api.DELETE("/questions/:id", s.apiClauseController.DeleteQuestion)
	}

	// // HTML routes group
	html := r.Group("/web")
	{
		html.GET("/iso_standards", s.webIsoStandardController.GetAllISOStandards)
		html.GET("/iso_standards/add", s.webIsoStandardController.RenderAddISOStandardForm)
		html.POST("/iso_standards/add", s.webIsoStandardController.CreateISOStandard)

		html.GET("/clauses", s.webClauseController.GetAllClauses)
		html.POST("/clauses/add", s.webClauseController.CreateClause)

		html.GET("/sections", s.webClauseController.GetAllSections)
		html.POST("/sections/add", s.webClauseController.CreateSection)

		html.GET("/questions", s.webClauseController.GetAllQuestions)
		html.POST("/questions/add", s.webClauseController.CreateQuestion)
	}

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
