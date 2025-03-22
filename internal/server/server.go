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
	webControllers "ISO_Auditing_Tool/pkg/controllers/web"
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/services"

	// webControllers "ISO_Auditing_Tool/cmd/web/controllers"
	"ISO_Auditing_Tool/pkg/repositories"
)

type Server struct {
	port                           int
	db                             database.Service
	eventBus                       *events.EventBus
	apiDraftController             *apiControllers.ApiDraftController
	webStandardController          *webControllers.WebStandardController
	apiMaterializedQueryController *apiControllers.ApiMaterializedQueryController

	// Add more Controllers here...
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()

	eventBus := events.NewEventBus()

	eventBus.Subscribe(events.MaterializedQueryCreated, events.LoggingHandler())
	eventBus.Subscribe(events.MaterializedQueryUpdated, events.LoggingHandler())
	eventBus.Subscribe(events.MaterializedQueryRefreshRequested, events.LoggingHandler())

	apiDraftRepo, _ := repositories.NewDraftRepository(db.DB())
	apiDraftService := services.NewDraftService(apiDraftRepo)
	apiDraftController := apiControllers.NewAPIDraftController(apiDraftService)

	apiMaterializedQueryRepo, _ := repositories.NewMaterializedQueriesRepository(db.DB())
	apiMaterializedQueryService := services.NewMaterializedQueryService(apiMaterializedQueryRepo, eventBus)
	apiMaterializedQueryController := apiControllers.NewApiMaterializedQueryController(apiMaterializedQueryService)

	webStandardRepo, _ := repositories.NewStandardRepository(db.DB())
	webStandardService := services.NewStandardService(webStandardRepo)
	webStandardController := webControllers.NewWebStandardController(webStandardService)

	// Add more repos, services, and controllers here...

	return &Server{
		port:                           port,
		db:                             db,
		eventBus:                       eventBus,
		apiDraftController:             apiDraftController,
		apiMaterializedQueryController: apiMaterializedQueryController,
		webStandardController:          webStandardController,

		// Add more Controllers here...
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
