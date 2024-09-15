package main

import (
	"livender-be/config"
	"livender-be/model"
	"livender-be/repository"
	"livender-be/rest"
	"log"

	"github.com/gin-gonic/gin"
)

func SetupRouter(genreRepo repository.GenreRepo, userRepo repository.UserRepo) *gin.Engine {
	e := gin.Default()

	// e.Use() // Cors
	// e.Use() // Authorization
	// e.Use() // Request-response logging

	rest.GenreRoutes(e, genreRepo)
	rest.UserRoutes(e, userRepo)

	return e
}

func main() {

	db := config.NewPostgresDB()
	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalln("Failed to connect to database.")
	}

	// dbConnection.Migrator().DropTable(&model.Book{}, &model.Genre{}, &model.Order{}, &model.User{})
	err = dbConnection.AutoMigrate(&model.Book{}, &model.Genre{}, &model.Order{}, &model.User{})
	if err != nil {
		log.Fatalln("Failed to auto migrate.")
	}

	genreRepo := repository.NewGenreRepo(dbConnection)
	userRepo := repository.NewUserRepo(dbConnection)

	r := SetupRouter(genreRepo, userRepo)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
