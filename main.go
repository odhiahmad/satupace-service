package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"run-sync/helper"
	"run-sync/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	helper.InitWhatsApp()

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "release"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
		log.Println("🛠 GIN_MODE=debug")
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Println("🚀 GIN_MODE=release")
	}

	r := routes.SetupRouter()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("✅ Server running on port %s...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server shutdown complete.")
}
