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
	http.Redirect(w, req, "http://"+req.Host+req.URL.String(), http.StatusFound)
}

func main() {

	mode := os.Getenv("GIN_MODE")
	r := routes.SetupRouter()

	if mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		go http.ListenAndServe(":80", http.HandlerFunc(redirectHTTP))
		err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", r)
		if err != nil {
			panic("Failed to start server: " + err.Error())
		}
	}

	r.Use(CORSMiddleware())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	r.Run(":" + port)
}
