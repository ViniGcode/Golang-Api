package controller

import (
	"gin-rest-api/api"
	"gin-rest-api/repository"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// struct vazio para agrupar os métodos do controller
type BookController struct{}

//List : Lista todos livros
func (bc *BookController) List(c *gin.Context) {
	repository := new(repository.BookRepository)
	var books = repository.List()
	c.JSON(200, books)
}

//Get : Obtém um livro pelo Id
func (bc *BookController) Get(c *gin.Context) {
	id := c.Param("id")
	repository := new(repository.BookRepository)
	book := repository.Get(id)
	c.JSON(http.StatusOK, book)
}

//Update : Atualiza um livro
func (bc *BookController) Update(c *gin.Context) {
	log.Println("Inside Update")
	var book api.Book
	repository := new(repository.BookRepository)

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		_, err := repository.Update(book)

		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"message": "Book record uodated succesfully !!!",
			})
		}
	}

}

//Delete : Deleta livro pelo Id
func (bc *BookController) Delete(c *gin.Context) {
	log.Println("Inside delete")
	id := c.Param("id")
	repository := new(repository.BookRepository)

	_, err := repository.Delete(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Book with ID '" + id + "' deleted succesfully !!!",
		})
	}

}

//Create : Cria um novo Livro.
func (bc *BookController) Create(c *gin.Context) {
	var book api.Book
	repository := new(repository.BookRepository)

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		bookID := repository.Create(book)
		//TODO : enviar loc header
		c.JSON(http.StatusOK, bookID)
	}
}
