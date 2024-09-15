package service

import (
	"errors"
	"livender-be/model"
	"livender-be/repository"
	"livender-be/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) UserService {
	return UserService{repo}
}

func (us UserService) Create(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	// Hash password (using bcrypt or any preferred hashing library)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password."})
		return
	}
	user.Password = string(hashedPassword)

	err = us.userRepo.Store(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully.", "user": user})
}

func (us UserService) GetAll(c *gin.Context) {
	var users []model.User
	err := us.userRepo.FindAll(&users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve users."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (us UserService) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User ID is invalid."})
		return
	}

	var user model.User

	err = us.userRepo.FindByID(id, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (us UserService) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User ID is invalid."})
		return
	}
	var user model.User

	// Fetch existing user
	err = us.userRepo.FindByID(id, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user."})
		return
	}

	// Bind updated data
	err = c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	// Optionally handle password update
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password."})
			return
		}
		user.Password = string(hashedPassword)
	}

	// Save updated user
	err = us.userRepo.Update(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully.", "user": user})
}

func (us UserService) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User ID is invalid."})
		return
	}

	// Fetch the user to confirm existence
	var user model.User
	err = us.userRepo.FindByID(id, &user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user."})
		return
	}

	// Delete the user
	err = us.userRepo.Delete(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully."})
}

func (us UserService) Register(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password."})
		return
	}
	user.Password = string(hashedPassword)

	err = us.userRepo.Store(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
}

func (us UserService) Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := c.ShouldBindJSON(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	var user model.User
	err = us.userRepo.FindByUsername(credentials.Username, &user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password."})
		return
	}

	// Generate JWT token
	token, err := util.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Example protected route to get user profile
func (us UserService) GetProfile(c *gin.Context) {
	userID := int(c.MustGet("userID").(uint))

	var user model.User
	err := us.userRepo.FindByID(userID, &user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
