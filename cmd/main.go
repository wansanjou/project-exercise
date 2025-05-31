package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wansanjou/backend-exercise-user-api/config"
	"github.com/wansanjou/backend-exercise-user-api/infrastructures"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/services"
	handlers "github.com/wansanjou/backend-exercise-user-api/internal/handlers/http"
	"github.com/wansanjou/backend-exercise-user-api/internal/repositories"
)

func main() {
	config.Init()
	db := infrastructures.NewMongoDB()
	r := gin.Default()

	ur := repositories.NewUserRepository(db, config.Get().Mongo.Database)
	us := services.NewUserService(ur)
	uh := handlers.NewUserHandler(us)

	as := services.NewAuthService(ur)
	ah := handlers.NewAuthHandler(as)

	api := r.Group("/api/v1")

	uh.UserRoutes(api)
	ah.AuthRoutes(api)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Background goroutine to count users
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping user count background job...")
				return
			case <-ticker.C:
				count, err := us.CountUsers(context.Background())
				if err != nil {
					log.Printf("failed to count users: %v", err)
					continue
				}
				log.Printf("Total users in DB: %d", count)
			}
		}
	}()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("HTTP server listening on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP serve failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Received shutdown signal...")

	// Start graceful shutdown
	cancel()

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := httpServer.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	wg.Wait()
	log.Println("Server gracefully stopped.")
}
