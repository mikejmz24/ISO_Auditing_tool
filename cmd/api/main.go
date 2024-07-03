// package main
//
// import (
//
//	"ISO_Auditing_Tool/cmd/server"
//	"context"
//	"log"
//	"net/http"
//	"os"
//	"os/signal"
//	"syscall"
//	"time"
//
//	"github.com/joho/godotenv"
//
// )
//
//	func main() {
//		// Load environment variables from .env file
//		if err := godotenv.Load(); err != nil {
//			log.Fatalf("Error loading .env file: %v", err)
//		}
//
//		srv := server.NewServer()
//		httpServer := srv.Start()
//
//		// Channel to listen for interrupt signals
//		stop := make(chan os.Signal, 1)
//		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
//
//		// Run server in a separate goroutine so that it doesn't block
//		go func() {
//			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//				log.Fatalf("could not listen on %s: %v\n", httpServer.Addr, err)
//			}
//		}()
//		log.Printf("Server is ready to handle requests at %s", httpServer.Addr)
//
//		// Block until we receive a signal
//		<-stop
//		log.Println("Server is shutting down...")
//
//		// Create a deadline to wait for.
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//
//		// Attempt a graceful shutdown
//		if err := httpServer.Shutdown(ctx); err != nil {
//			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
//		}
//		log.Println("Server stopped")
//	}
package main

import (
	"ISO_Auditing_Tool/cmd/internal/database"
	"ISO_Auditing_Tool/cmd/server"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check for command line arguments for migration and seeding
	if len(os.Args) > 1 {
		dbService := database.New()
		defer dbService.Close()

		switch os.Args[1] {
		case "migrate":
			if err := dbService.Migrate(); err != nil {
				log.Fatalf("Failed to migrate database: %v", err)
			}
			log.Println("Database migrated successfully")
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
		default:
			log.Fatalf("Unknown command: %s", os.Args[1])
		}
		return
	}

	srv := server.NewServer()
	httpServer := srv.Start()

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
