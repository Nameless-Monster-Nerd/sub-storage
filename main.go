package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nameless-Monster-Nerd/subtitle/src/routes/proxy"
	"github.com/nameless-Monster-Nerd/subtitle/src/routes/subs"
)

func main() {
	router := gin.Default()

	// Custom CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Handle preflight (OPTIONS) requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/subs/:id", sub.Sub)
	router.GET("/proxy.vtt", proxy.Proxy)

	router.Run(":8080") // http://localhost:8080/subs/123?ss=1&ep=2
}
