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
	port                 int
	db                   database.Service
	apiClauseController  *apiControllers.ApiClauseController
	htmlClauseController *webControllers.HtmlClauseController
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()

	// Call Migrate method to create tables
	if err := db.Migrate(); err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
		os.Exit(1)
	}

	clauseRepo := repositories.NewClauseRepository(db.DB())
	apiClauseController := apiControllers.NewApiClauseController(clauseRepo)
	htmlClauseController := webControllers.NewHtmlClauseController(clauseRepo)

	NewServer := &Server{
		port:                 port,
		db:                   db,
		apiClauseController:  apiClauseController,
		htmlClauseController: htmlClauseController,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(db.DB()),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
