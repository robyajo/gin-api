package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service berjalan normal",
		})
	})

	fmt.Println("Server berjalan di http://localhost:8080")
	r.Run(":8080")
}
