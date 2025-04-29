package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"ISO_Auditing_Tool/internal/database"
	apiControllers "ISO_Auditing_Tool/pkg/controllers/api"
	webControllers "ISO_Auditing_Tool/pkg/controllers/web"
	"ISO_Auditing_Tool/pkg/events"
	"ISO_Auditing_Tool/pkg/repositories"
	"ISO_Auditing_Tool/pkg/services"
)

// Config holds all configuration for the server
type Config struct {
	Port           int              `json:"port"`
	Host           string           `json:"host"`
	ReadTimeout    time.Duration    `json:"read_timeout"`
	WriteTimeout   time.Duration    `json:"write_timeout"`
	IdleTimeout    time.Duration    `json:"idle_timeout"`
	DatabaseConfig *database.Config `json:"database_config"`
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() (*Config, error) {
	portStr := os.Getenv("PORT")
	port := 8080 // Default port

	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT value: %w", err)
		}
	}

	// Read timeout with default
	readTimeoutStr := os.Getenv("READ_TIMEOUT")
	readTimeout := 10 * time.Second
	if readTimeoutStr != "" {
		readTimeoutSec, err := strconv.Atoi(readTimeoutStr)
		if err == nil {
			readTimeout = time.Duration(readTimeoutSec) * time.Second
		}
	}

	// Write timeout with default
	writeTimeoutStr := os.Getenv("WRITE_TIMEOUT")
	writeTimeout := 30 * time.Second
	if writeTimeoutStr != "" {
		writeTimeoutSec, err := strconv.Atoi(writeTimeoutStr)
		if err == nil {
			writeTimeout = time.Duration(writeTimeoutSec) * time.Second
		}
	}

	// Idle timeout with default
	idleTimeoutStr := os.Getenv("IDLE_TIMEOUT")
	idleTimeout := time.Minute
	if idleTimeoutStr != "" {
		idleTimeoutSec, err := strconv.Atoi(idleTimeoutStr)
		if err == nil {
			idleTimeout = time.Duration(idleTimeoutSec) * time.Second
		}
	}

	// Load database configuration
	dbConfig := database.LoadConfigFromEnv()

	return &Config{
		Port:           port,
		Host:           os.Getenv("HOST"),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
		DatabaseConfig: dbConfig,
	}, nil
}

type Server struct {
	config                             *Config
	db                                 database.Service
	eventBus                           *events.EventBus
	apiDraftController                 *apiControllers.ApiDraftController
	webStandardController              *webControllers.WebStandardController
	apiMaterializedJSONQueryController *apiControllers.ApiMaterializedJSONQueryController
}

// NewServer creates a new server instance with the given configuration
func NewServer() (*Server, error) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize database with config
	db := database.NewWithConfig(config.DatabaseConfig)

	// Create event bus with error handling
	eventBus := events.NewEventBus()

	// Setup event subscribers with error handling callbacks
	eventBus.Subscribe(events.MaterializedQueryCreated, events.LoggingHandler())
	eventBus.Subscribe(events.MaterializedQueryUpdated, events.LoggingHandler())
	eventBus.Subscribe(events.MaterializedQueryRefreshRequested, events.LoggingHandler())

	// Setup repositories
	draftRepo, err := repositories.NewDraftRepository(db.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to create draft repository: %w", err)
	}

	materializedJSONQueryRepo, err := repositories.NewMaterializedQueriesJSONRepository(db.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to create materialized JSON query repository: %w", err)
	}

	materializedHTMLQueryRepo, err := repositories.NewMaterializedQueriesHTMLRepository(db.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to create materialized HTML query repository: %w", err)
	}

	standardRepo, err := repositories.NewStandardRepository(db.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to create standard repository: %w", err)
	}

	requirementRepo, err := repositories.NewRequirementRepository(db.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to create requirement repository: %w", err)
	}

	questionRepo, err := repositories.NewQuestionRepository(db.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to create question repository: %w", err)
	}

	evidenceRepo, err := repositories.NewEvidenceRepository(db.DB())
	if err != nil {
		return nil, fmt.Errorf("failed to create evidence repository: %w", err)
	}

	// Setup services
	draftService := services.NewDraftService(draftRepo)
	// apiMaterializedQueryService := services.NewMaterializedJSONService(apiMaterializedQueryRepo, eventBus)
	materializedJSONQueryService := services.NewMaterializedJSONService(materializedJSONQueryRepo, standardRepo, requirementRepo, questionRepo, evidenceRepo, eventBus)
	htmlCacheService := services.NewHTMLCacheService(materializedHTMLQueryRepo, materializedJSONQueryRepo, standardRepo, requirementRepo, eventBus)
	standardService := services.NewStandardService(standardRepo)

	// Setup controllers
	apiDraftController := apiControllers.NewAPIDraftController(draftService)
	apiMaterializedQueryController := apiControllers.NewApiMaterializedJSONQueryController(materializedJSONQueryService, htmlCacheService, eventBus)
	webStandardController := webControllers.NewWebStandardController(standardService)

	return &Server{
		config:                             config,
		db:                                 db,
		eventBus:                           eventBus,
		apiDraftController:                 apiDraftController,
		apiMaterializedJSONQueryController: apiMaterializedQueryController,
		webStandardController:              webStandardController,
	}, nil
}

// Start initializes and starts the HTTP server
func (s *Server) Start() (*http.Server, error) {
	// Ensure database is ready
	if err := s.db.Ping(); err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	// Build server address
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	if s.config.Host == "" {
		addr = fmt.Sprintf(":%d", s.config.Port)
	}

	// Declare server config
	server := &http.Server{
		Addr:         addr,
		Handler:      s.RegisterRoutes(s.db.DB()),
		IdleTimeout:  s.config.IdleTimeout,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	// Log server startup
	log.Printf("Starting server on %s", addr)
	return server, nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	// Close database connections
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("error closing database connections: %w", err)
	}

	return nil
}
