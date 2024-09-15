package rest

import (
	"livender-be/repository"
	"livender-be/service"

	"github.com/gin-gonic/gin"
)

func GenreRoutes(r *gin.Engine, genreRepo repository.GenreRepo) {
	s := service.NewGenreService(genreRepo)

	g := r.Group("/genres")
	g.GET("", s.GetAll)
	g.GET(":id", s.GetByID)
	g.POST("", s.Create)

	b := r.Group("/books")
	b.GET("/genre/:id", s.GetBooksByGenre) // get books by genre
}
