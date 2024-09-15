package main

import (
	"livender-be/config"
	"livender-be/model"
	"livender-be/repository"
	"livender-be/rest"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(genreRepo repository.GenreRepo, userRepo repository.UserRepo, bookRepo repository.BookRepo) *gin.Engine {
	e := gin.Default()

	// e.Use() // Cors
	// e.Use() // Authorization
	// e.Use() // Request-response logging

	rest.GenreRoutes(e, genreRepo)
	rest.UserRoutes(e, userRepo)
	rest.BookRoutes(e, bookRepo)

	return e
}

func truncateTables(db *gorm.DB) error {
	// Truncate tables with cascade
	err := db.Exec("TRUNCATE TABLE books, genres, orders, users RESTART IDENTITY CASCADE").Error
	if err != nil {
		return err
	}
	return nil
}

func main() {

	db := config.NewPostgresDB()
	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalln("Failed to connect to database.")
	}

	// dbConnection.Migrator().DropTable(&model.Book{}, &model.Genre{}, &model.Order{}, &model.User{})
	// Truncate tables
	// if err := truncateTables(dbConnection); err != nil {
	// 	log.Fatalln("Failed to truncate tables:", err)
	// }

	err = dbConnection.AutoMigrate(&model.User{}, &model.Book{}, &model.Genre{}, &model.Order{})
	if err != nil {
		log.Fatalln("Failed to auto migrate.")
	}

	genreFantasy := model.Genre{Name: "Fantasy"}
	genreFiction := model.Genre{Name: "Fiction"}
	genreNonFiction := model.Genre{Name: "Non-Fiction"}
	genreScienceFiction := model.Genre{Name: "Adventure"}
	genreMystery := model.Genre{Name: "Mystery"}
	genreHorror := model.Genre{Name: "Horror"}

	dbConnection.Create(&genreFantasy)
	dbConnection.Create(&genreFiction)
	dbConnection.Create(&genreNonFiction)
	dbConnection.Create(&genreScienceFiction)
	dbConnection.Create(&genreMystery)
	dbConnection.Create(&genreHorror)

	genreRepo := repository.NewGenreRepo(dbConnection)
	userRepo := repository.NewUserRepo(dbConnection)
	bookRepo := repository.NewBookRepo(dbConnection)

	r := SetupRouter(genreRepo, userRepo, bookRepo)

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
