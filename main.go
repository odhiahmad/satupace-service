package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/routes"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func redirectHTTP(w http.ResponseWriter, req *http.Request) {
	// Ubah ke HTTPS (bukan HTTP) untuk redirect production
	http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusFound)
}

func main() {
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "release"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin mode SEBELUM router dibuat
	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := routes.SetupRouter()
	r.Use(CORSMiddleware())

	if mode == "debug" {
		// Development mode: pakai HTTP di port custom
		if err := r.Run(":" + port); err != nil {
			panic("Failed to start debug server: " + err.Error())
		}
	} else {
		// Production mode:
		// Redirect HTTP (port 80) ke HTTPS (port 443)
		go func() {
			if err := http.ListenAndServe(":80", http.HandlerFunc(redirectHTTP)); err != nil {
				panic("Failed to start redirect server: " + err.Error())
			}
		}()

		// Jalankan server HTTPS (port 443)
		if err := r.RunTLS(":443", "cert.pem", "key.pem"); err != nil {
			panic("Failed to start TLS server: " + err.Error())
		}
	}
}
