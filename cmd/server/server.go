package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	apiControllers "ISO_Auditing_Tool/cmd/api/controllers"
	"ISO_Auditing_Tool/cmd/internal/database"
	webControllers "ISO_Auditing_Tool/cmd/web/controllers"
	"ISO_Auditing_Tool/pkg/repositories"
)

type Server struct {
	port                      int
	db                        database.Service
	apiIsoStandardController  *apiControllers.ApiIsoStandardController
	apiClauseController       *apiControllers.ApiClauseController
	htmlIsoStandardController *webControllers.HtmlIsoStandardController
	htmlClauseController      *webControllers.HtmlClauseController
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()

	// Call Migrate method to create tables
	if err := db.Migrate(); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	clauseRepo := repositories.NewClauseRepository(db.DB())
	apiStandardRepo := repositories.NewIsoStandardRepository(db.DB())

	apiClauseController := apiControllers.NewApiClauseController(clauseRepo)
	apiIsoStandardController := apiControllers.NewApiIsoStandardController(apiStandardRepo)
	htmlClauseController := webControllers.NewHtmlClauseController(clauseRepo)
	htmlIsoStandardController := webControllers.NewHtmlIsoStandardController(apiStandardRepo)

	return &Server{
		port:                      port,
		db:                        db,
		apiIsoStandardController:  apiIsoStandardController,
		apiClauseController:       apiClauseController,
		htmlIsoStandardController: htmlIsoStandardController,
		htmlClauseController:      htmlClauseController,
	}
}

func (s *Server) Start() *http.Server {
	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(s.db.DB()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
