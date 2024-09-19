package main

import (
	"livender-be/config"
	"livender-be/model"
	"livender-be/repository"
	"livender-be/rest"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(genreRepo repository.GenreRepo, userRepo repository.UserRepo, bookRepo repository.BookRepo, orderRepo repository.OrderRepo) *gin.Engine {
	e := gin.Default()

	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	rest.GenreRoutes(e, genreRepo)
	rest.UserRoutes(e, userRepo)
	rest.BookRoutes(e, bookRepo)
	rest.OrderRoutes(e, orderRepo, bookRepo)

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

	// if err := truncateTables(dbConnection); err != nil {
	// 	log.Fatalln("Failed to truncate tables:", err)
	// }

	err = dbConnection.AutoMigrate(&model.User{}, &model.Book{}, &model.Genre{}, &model.Order{})
	if err != nil {
		log.Fatalln("Failed to auto migrate.")
	}

	var genres []model.Genre
	err = dbConnection.Find(&genres).Error
	if err != nil {
		log.Fatalln("Failed to find genres:", err)
	}
	if len(genres) == 0 {
		genres = []model.Genre{
			{Name: "Fantasy"},
			{Name: "Fiction"},
			{Name: "Non-Fiction"},
			{Name: "Adventure"},
			{Name: "Mystery"},
			{Name: "Horror"},
		}
		err = dbConnection.Create(&genres).Error
		if err != nil {
			log.Fatalln("Failed to create genres:", err)
		}
	}

	genreRepo := repository.NewGenreRepo(dbConnection)
	userRepo := repository.NewUserRepo(dbConnection)
	bookRepo := repository.NewBookRepo(dbConnection)
	orderRepo := repository.NewOrderRepo(dbConnection)

	r := SetupRouter(genreRepo, userRepo, bookRepo, orderRepo)

	r.Run(":8080")
}
