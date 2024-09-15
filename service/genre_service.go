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

type GenreService struct {
	genreRepo repository.GenreRepo
}

func NewGenreService(repo repository.GenreRepo) GenreService {
	return GenreService{repo}
}

func (gs GenreService) Create(c *gin.Context) {
	var g model.Genre
	err := c.ShouldBindJSON(&g)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	err = gs.genreRepo.Store(&g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create Genre."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Genre has been created successfully.", "genre": g})

}

func (gs GenreService) GetAll(c *gin.Context) {
	var genres []model.Genre
	err := gs.genreRepo.FindAll(&genres)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch Genres."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"genres": genres})
}

func (gs GenreService) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Genre ID is invalid."})
		return
	}

	var genre model.Genre

	err = gs.genreRepo.FindByID(id, &genre)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Genre not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve the genre."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"genre": genre})
}

func (gs GenreService) GetBooksByGenre(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Genre ID is invalid."})
		return
	}

	var genre model.Genre

	err = gs.genreRepo.FindByID(id, &genre)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Genre not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve the genre."})
		return
	}

	var books []model.Book
	err = gs.genreRepo.FindBooksByGenre(id, &books)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"books": books})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}
