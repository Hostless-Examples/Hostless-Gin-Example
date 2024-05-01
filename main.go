package main

import (
	"fmt"
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
)


func main() {

	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html" ,gin.H{
			"message": "success",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router.Run(fmt.Sprintf(":%s", port)) // listen and serve on 0.0.0.0:8080
}