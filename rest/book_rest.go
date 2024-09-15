package rest

import (
	"livender-be/middleware"
	"livender-be/repository"
	"livender-be/service"

	"github.com/gin-gonic/gin"
)

func BookRoutes(r *gin.Engine, bookRepo repository.BookRepo) {
	bs := service.NewBookService(bookRepo)

	// Public routes for books (e.g., for viewing all books)
	b := r.Group("/books")
	b.GET("", bs.GetAll)                    // Get all books
	b.GET("/:id", bs.GetByID)               // Get book by ID
	b.GET("/user/:id", bs.GetBooksByUserID) // get books owned by user

	// Protected routes for book creation, update, delete (requires JWT)
	protected := r.Group("/books")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("", bs.Create)       // Create a new book
	protected.PUT("/:id", bs.Update)    // Update book by ID
	protected.DELETE("/:id", bs.Delete) // Delete book by ID
}
