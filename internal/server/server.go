package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	// "ISO_Auditing_Tool/cmd/api/controllers"
	"ISO_Auditing_Tool/cmd/api/controllers/api"
	"ISO_Auditing_Tool/cmd/api/controllers/html"
	"ISO_Auditing_Tool/cmd/api/repositories"
	"ISO_Auditing_Tool/internal/database"
)

type Server struct {
	port int
	db   database.Service
	// clauseController *controllers.ClauseController
	apiClauseController  *api.ApiClauseController
	htmlClauseController *html.HtmlClauseController
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	clauseRepo := repositories.NewClauseRepository(db.DB())
	// clauseController := controllers.NewClauseController(clauseRepo)
	apiClauseController := api.NewApiClauseController(clauseRepo)
	htmlClauseController := html.NewHtmlClauseController(clauseRepo)
	NewServer := &Server{
		port: port,
		db:   database.New(),
		// clauseController: clauseController,
		apiClauseController:  apiClauseController,
		htmlClauseController: htmlClauseController,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(&sql.DB{}),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
