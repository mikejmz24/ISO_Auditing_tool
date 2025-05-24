package server

import (
	"net/http"

	"ISO_Auditing_Tool/pkg/middleware"
	"github.com/gin-gonic/gin"

	"ISO_Auditing_Tool/internal/database"
	"database/sql"
)

func (s *Server) RegisterRoutes(db *sql.DB) http.Handler {
	r := gin.Default()
	s.db = database.New()

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	// r.Static("/assets", "../cmd/web/assets")
	r.Static("/web/assets", "cmd/web/assets")

	// API routes group
	api := r.Group("/api")
	api.Use(middleware.ErrorHandler())
	{
		api.POST("/drafts", s.apiDraftController.Create)
		api.PUT("/drafts/:id", s.apiDraftController.Update)
		api.GET("/drafts", s.apiDraftController.GetAll)
		// api.GET("/iso_standards", s.apiIsoStandardController.GetAllISOStandards)
		// api.GET("/iso_standards/:id", s.apiIsoStandardController.GetISOStandardByID)
		// api.POST("/iso_standards", s.apiIsoStandardController.CreateISOStandard)
		// api.PUT("/iso_standards/:id", s.apiIsoStandardController.UpdateISOStandard)
		// api.DELETE("/iso_standards/:id", s.apiIsoStandardController.DeleteISOStandard)
		api.GET("/query/:name", s.apiMaterializedJSONQueryController.GetByName)
		api.POST("/query", s.apiMaterializedJSONQueryController.CreateOrUpdateJSONQuery)
	}

	// // HTML routes group
	html := r.Group("/web")
	html.Use(middleware.ErrorHandler())
	{
		// html.GET("/iso_standards", s.webIsoStandardController.GetAllISOStandards)
		// html.GET("/iso_standards/add", s.webIsoStandardController.RenderAddISOStandardForm)
		// html.POST("/iso_standards", s.webIsoStandardController.CreateISOStandard)
		// html.GET("/iso_standards/:id", s.webIsoStandardController.GetISOStandardByID)
		html.GET("/standards/:id", s.webStandardController.GetByID)
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
