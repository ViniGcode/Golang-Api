package main

import (
	"gin-rest-api/controller"
	"log"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

//CustomMiddleware autentica todo request
func CustomMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Inside customMiddleware")
		c.Next()
	}
}

func main() {

	router := gin.Default()

	router.Use(CustomMiddleware())
	router.Use(cors.Default())

	api := router.Group("api/v1")
	{
		book := new(controller.BookController)
		api.GET("/books", book.List)
		api.GET("/books/:id", book.Get)
		api.POST("/books", book.Create)
		api.PUT("/books", book.Update)
		api.DELETE("/books/:id", book.Delete)
	}

	router.Run()
}
