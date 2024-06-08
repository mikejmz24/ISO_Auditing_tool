package server

import (
	"database/sql"
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
	clauseRepo := repositories.NewClauseRepository(db.DB())
	apiClauseController := apiControllers.NewApiClauseController(clauseRepo)
	htmlClauseController := webControllers.NewHtmlClauseController(clauseRepo)
	NewServer := &Server{
		port:                 port,
		db:                   database.New(),
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
