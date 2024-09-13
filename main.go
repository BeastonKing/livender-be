package main

import (
	"livender-be/config"
	"livender-be/controller"
	"livender-be/model"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	db := config.NewPostgresDB()
	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalln("Failed to connect to database.")
	}

	err = dbConnection.AutoMigrate(&model.Book{}, &model.Genre{}, &model.Order{}, &model.User{})
	if err != nil {
		log.Fatalln("Failed to auto migrate.")
	}

	r := gin.Default()
	r.POST("/genre", controller.CreateGenre)
	r.Run() // listen and serve on 0.0.0.0:8080
}
