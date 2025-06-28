package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/odhiahmad/kasirku-service/routes"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func redirectHTTP(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
}

func main() {
	// Load .env jika ada
	_ = godotenv.Load(".env")

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "release"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin mode
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup router
	r := routes.SetupRouter()
	r.Use(CORSMiddleware())

	// Debug mode (HTTP)
	if mode == "debug" {
		server := &http.Server{
			Addr:    ":" + port,
			Handler: r,
		}

		go func() {
			log.Printf("🚀 Server running in DEBUG mode on port %s\n", port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("❌ Failed to start debug server: %v", err)
			}
		}()

		// Wait for interrupt (graceful shutdown)
		gracefulShutdown(server)
		return
	}

	// Release mode (HTTPS)
	go func() {
		// Redirect HTTP to HTTPS
		log.Println("🌐 Redirect HTTP (80) → HTTPS (443)")
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectHTTP)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start redirect server: %v", err)
		}
	}()

	// Jalankan server HTTPS
	server := &http.Server{
		Addr:    ":443",
		Handler: r,
	}

	go func() {
		log.Println("🔒 Running HTTPS server on port 443")
		if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start TLS server: %v", err)
		}
	}()

	// Graceful shutdown untuk release juga
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
