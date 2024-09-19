package main_test

import (
	"encoding/json"
	"livender-be/config"
	"livender-be/model"
	"livender-be/repository"
	"livender-be/rest"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

func TestRegisterUser(t *testing.T) {
	db := config.NewPostgresDB()
	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalln("Failed to connect to database.")
	}

	tx := dbConnection.Begin()

	defer tx.Rollback()

	genreRepo := repository.NewGenreRepo(tx)
	userRepo := repository.NewUserRepo(tx)
	bookRepo := repository.NewBookRepo(tx)
	orderRepo := repository.NewOrderRepo(tx)

	r := SetupRouter(genreRepo, userRepo, bookRepo, orderRepo)
	user := model.User{
		Username: "testuser",
		Password: "testpassword",
		Fullname: "Test User",
	}

	w := httptest.NewRecorder()
	userJson, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users/register", strings.NewReader(string(userJson)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}
}

func TestLoginUser(t *testing.T) {
	db := config.NewPostgresDB()
	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalln("Failed to connect to database.")
	}

	tx := dbConnection.Begin()

	defer tx.Rollback()

	genreRepo := repository.NewGenreRepo(tx)
	userRepo := repository.NewUserRepo(tx)
	bookRepo := repository.NewBookRepo(tx)
	orderRepo := repository.NewOrderRepo(tx)

	r := SetupRouter(genreRepo, userRepo, bookRepo, orderRepo)
	user := model.User{
		Username: "testuser",
		Password: "testpassword",
		Fullname: "Test User",
	}

	w := httptest.NewRecorder()
	userJson, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users/register", strings.NewReader(string(userJson)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	user = model.User{
		Username: "testuser",
		Password: "testpassword",
	}

	w = httptest.NewRecorder()
	userJson, _ = json.Marshal(user)
	req, _ = http.NewRequest("POST", "/users/login", strings.NewReader(string(userJson)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	if err != nil {
		t.Errorf("Failed to login user: %v", err)
	}
}
