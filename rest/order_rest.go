package rest

import (
	"livender-be/middleware"
	"livender-be/repository"
	"livender-be/service"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine, orderRepo repository.OrderRepo, bookRepo repository.BookRepo) {
	os := service.NewOrderService(orderRepo, bookRepo)

	// Protected routes for creating orders (purchases)
	protected := r.Group("/orders")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/:id", os.GetByID)        // Purchase a book
	protected.POST("/purchase", os.Purchase) // Purchase a book

	// View user orders
	protected.GET("/user/:id", os.GetUserOrders) // Get all orders by user
}
