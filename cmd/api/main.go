package main

import (
	"ISO_Auditing_Tool/internal/database"
	"ISO_Auditing_Tool/internal/server"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	env := os.Getenv("ENV")
	envFile := ".env"
	if env == "production" {
		envFile = ".env.production"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading %s file: %v", envFile, err)
	}
}

func main() {
	// Load environment variables
	loadEnv()
	gin.SetMode(gin.DebugMode)

	// Check for command line arguments for migration and seeding
	if len(os.Args) > 1 {
		dbService := database.New()
		defer dbService.Close()

		switch os.Args[1] {
		case "migrate":
			var direction string
			var file string

			if len(os.Args) > 2 {
				direction = validateDirection(os.Args[2])

				if os.Args[3] != "" {
					file = os.Args[3]
				}
			}

			if err := dbService.Migrate(file, direction); err != nil {
				log.Fatalf("Failed to migrate database: %v", err)
			}
			log.Printf("Database migrated successfully (%s)", direction)

		case "seed":
			if err := dbService.Seed(); err != nil {
				log.Fatalf("Failed to seed database: %v", err)
			}
			log.Println("Database seeded successfully")

		case "truncate":
			if err := dbService.Truncate(); err != nil {
				log.Fatalf("Failed to truncate seeded data: %v", err)
			}
			log.Println("Seeded data truncated successfully")

		case "refresh":
			if err := dbService.RefreshDatabase(); err != nil {
				log.Fatalf("Failed to refresh database: %v", err)
			}
			log.Println("Refreshed database successfully")

		default:
			log.Fatalf("Unknown command: %s", os.Args[1])
		}
		return
	}

	srv, _ := server.NewServer()
	httpServer, _ := srv.Start()

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run server in a separate goroutine so that it doesn't block
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", httpServer.Addr, err)
		}
	}()
	log.Printf("Server is ready to handle requests at %s", httpServer.Addr)

	// Block until we receive a signal
	<-stop
	log.Println("Server is shutting down...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt a graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	log.Println("Server stopped")
}

func validateDirection(input string) string {
	if input == "up" || input == "down" {
		return input
	}
	return "up"
}
