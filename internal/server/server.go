package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"ISO_Auditing_Tool/internal/database"
	apiControllers "ISO_Auditing_Tool/pkg/controllers/api"
	"ISO_Auditing_Tool/pkg/services/api"

	// webControllers "ISO_Auditing_Tool/cmd/web/controllers"
	"ISO_Auditing_Tool/pkg/repositories"
)

type Server struct {
	port               int
	db                 database.Service
	apiDraftController *apiControllers.ApiDraftController
	// apiIsoStandardController *apiControllers.ApiIsoStandardController
	// apiClauseController      *apiControllers.ApiClauseController
	// webIsoStandardController *webControllers.WebIsoStandardController
	// webClauseController      *webControllers.HtmlClauseController
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()

	// if err := db.Migrate(); err != nil {
	// 	fmt.Printf("Failed to migrate database: %v\n", err)
	// 	os.Exit(1)
	// }
	//
	// if err := db.Seed(); err != nil {
	// 	fmt.Printf("Failed to seed the database: %v\n", err)
	// 	os.Exit(1)
	// }

	// clauseRepo := repositories.NewClauseRepository(db.DB())
	// apiStandardRepo := repositories.NewIsoStandardRepository(db.DB())
	apiDraftRepo, _ := repositories.NewDraftRepository(db.DB())
	apiDraftService := services.NewDraftService(apiDraftRepo)
	apiDraftController := apiControllers.NewAPIDraftController(apiDraftService)

	// apiClauseController := apiControllers.NewApiClauseController(clauseRepo)
	// apiIsoStandardController := apiControllers.NewApiIsoStandardController(apiStandardRepo)
	// htmlClauseController := webControllers.NewHtmlClauseController(clauseRepo)
	// htmlIsoStandardController := webControllers.NewWebIsoStandardController(apiIsoStandardController)

	return &Server{
		port:               port,
		db:                 db,
		apiDraftController: apiDraftController,
		// 	apiIsoStandardController: apiIsoStandardController,
		// 	apiClauseController:      apiClauseController,
		// 	webIsoStandardController: htmlIsoStandardController,
		// 	webClauseController:      htmlClauseController,
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
