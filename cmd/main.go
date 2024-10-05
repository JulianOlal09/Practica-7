package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	fmt.Println("Running App")
	//Tomar archivos de la carpeta templates
	router.LoadHTMLGlob("templates/*")
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	router.Run(":8000")
}
