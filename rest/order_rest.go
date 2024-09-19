package rest

import (
	"livender-be/middleware"
	"livender-be/repository"
	"livender-be/service"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.Engine, orderRepo repository.OrderRepo, bookRepo repository.BookRepo) {
	os := service.NewOrderService(orderRepo, bookRepo)

	protected := r.Group("/orders")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/:id", os.GetByID)
	protected.POST("/purchase", os.Purchase)

	protected.GET("/user/:id", os.GetUserOrders)
}
