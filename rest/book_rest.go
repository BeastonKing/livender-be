package rest

import (
	"livender-be/middleware"
	"livender-be/repository"
	"livender-be/service"

	"github.com/gin-gonic/gin"
)

func BookRoutes(r *gin.Engine, bookRepo repository.BookRepo) {
	bs := service.NewBookService(bookRepo)

	b := r.Group("/books")
	b.GET("", bs.GetAll)
	b.GET("/:id", bs.GetByID)
	b.GET("/user/:id", bs.GetBooksByUserID)

	protected := r.Group("/books")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("", bs.Create)
	protected.PUT("/:id", bs.Update)
	protected.DELETE("/:id", bs.Delete)
}
