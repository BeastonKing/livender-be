package service

import (
	"errors"
	"livender-be/model"
	"livender-be/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookService struct {
	bookRepo repository.BookRepo
}

func NewBookService(repo repository.BookRepo) BookService {
	return BookService{repo}
}

// Create a new book
func (bs BookService) Create(c *gin.Context) {
	var book model.Book
	err := c.ShouldBindJSON(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	err = bs.bookRepo.Store(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create book."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book created successfully.", "book": book})
}

// Get all books
func (bs BookService) GetAll(c *gin.Context) {
	var books []model.Book
	err := bs.bookRepo.FindAll(&books)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch books."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

// Get a book by ID
func (bs BookService) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book ID is invalid."})
		return
	}

	var book model.Book
	err = bs.bookRepo.FindByID(id, &book)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book})
}

// Update a book by ID
// func (bs BookService) Update(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Book ID is invalid."})
// 		return
// 	}

// 	var book model.Book
// 	err = bs.bookRepo.FindByID(id, &book)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
// 		return
// 	}

// 	err = c.ShouldBindJSON(&book)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
// 		return
// 	}

// 	book.ID = uint(id) // Ensure we don't modify the ID in the request body
// 	err = bs.bookRepo.Update(&book)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update book."})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully.", "book": book})
// }

func (bs BookService) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book ID is invalid."})
		return
	}

	var updatedBook model.Book
	err = c.ShouldBindJSON(&updatedBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	var book model.Book
	err = bs.bookRepo.FindByID(id, &book)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	// Update the book details
	book.Title = updatedBook.Title
	book.Author = updatedBook.Author
	book.ReleaseYear = updatedBook.ReleaseYear
	book.Age = updatedBook.Age
	book.UserID = updatedBook.UserID
	book.IsSold = updatedBook.IsSold

	// Clear existing genres and add the new ones
	err = bs.bookRepo.ClearGenres(&book) // Clears the many-to-many relationship
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to clear genres."})
		return
	}

	// Assign new genres
	book.Genres = updatedBook.Genres

	err = bs.bookRepo.Update(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update book."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully.", "book": book})
}

// Delete a book
func (bs BookService) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book ID is invalid."})
		return
	}

	var book model.Book
	err = bs.bookRepo.FindByID(id, &book)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve book."})
		return
	}

	err = bs.bookRepo.Delete(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete book."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully."})
}

func (bs BookService) GetBooksByUserID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is invalid."})
		return
	}

	var books []model.Book
	err = bs.bookRepo.FindAllBooksOwnedByUser(id, &books)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"books": books})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}
