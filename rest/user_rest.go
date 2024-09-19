package rest

import (
	"livender-be/middleware"
	"livender-be/repository"
	"livender-be/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userRepo repository.UserRepo) {
	us := service.NewUserService(userRepo)

	u := r.Group("/users")
	u.POST("/register", us.Register)
	u.POST("/login", us.Login)

	protected := r.Group("/users")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("", us.GetAll)
	protected.GET("/:id", us.GetByID)
	protected.PUT("/:id", us.Update)
	protected.DELETE("/:id", us.Delete)

	r.Group("/profile").Use(middleware.AuthMiddleware()).GET("", us.GetProfile)
}
