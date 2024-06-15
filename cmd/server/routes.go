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
	{
		api.GET("/clauses", s.apiClauseController.GetAllClauses)
		api.GET("/clauses/:id", s.apiClauseController.GetClauseByID)
		api.POST("/clauses", s.apiClauseController.CreateClause)
		api.PUT("/clauses/:id", s.apiClauseController.UpdateClause)
		api.DELETE("/clauses/:id", s.apiClauseController.DeleteClause)

		api.GET("/iso_standards", s.apiIsoStandardController.GetAllISOStandards)
		api.GET("/iso_standards/:id", s.apiIsoStandardController.GetISOStandardByID)
		api.POST("/iso_standards", s.apiIsoStandardController.CreateISOStandard)
		api.PUT("/iso_standards/:id", s.apiIsoStandardController.UpdateISOStandard)
		api.DELETE("/iso_standards/:id", s.apiIsoStandardController.DeleteISOStandard)

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
	html := r.Group("/html")
	{
		html.GET("/clauses", s.htmlClauseController.GetAllClauses)
		// html.GET("/clauses/add", func(c *gin.Context) {
		// 	templ.Handler(templates.AddClause()).ServeHTTP(c.Writer, c.Request)
		// })
		html.POST("/clauses/add", s.htmlClauseController.CreateClause)

		html.GET("/iso_standards", s.htmlIsoStandardController.GetAllISOStandards)
		html.GET("/iso_standards/add", s.htmlIsoStandardController.RenderAddISOStandardForm)
		html.POST("/iso_standards/add", s.htmlIsoStandardController.CreateISOStandard)

		html.GET("/sections", s.htmlClauseController.GetAllSections)
		// html.GET("/sections/add", func(c *gin.Context) {
		// 	templ.Handler(templates.AddSection()).ServeHTTP(c.Writer, c.Request)
		// })
		html.POST("/sections/add", s.htmlClauseController.CreateSection)

		html.GET("/questions", s.htmlClauseController.GetAllQuestions)
		// html.GET("/questions/add", func(c *gin.Context) {
		// 	templ.Handler(templates.AddQuestion()).ServeHTTP(c.Writer, c.Request)
		// })
		html.POST("/questions/add", s.htmlClauseController.CreateQuestion)
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
