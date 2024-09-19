package service

import (
	"errors"
	"livender-be/model"
	"livender-be/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo repository.OrderRepo
	bookRepo  repository.BookRepo
}

func NewOrderService(orderRepo repository.OrderRepo, bookRepo repository.BookRepo) OrderService {
	return OrderService{orderRepo, bookRepo}
}

func (os OrderService) Purchase(c *gin.Context) {
	var req struct {
		BookID uint `json:"book_id"`
		UserID uint `json:"user_id"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request."})
		return
	}

	var book model.Book
	err = os.bookRepo.FindByID(int(req.BookID), &book)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve book."})
		return
	}

	if book.IsSold {
		c.JSON(http.StatusBadRequest, gin.H{"message": "This book is already sold."})
		return
	}

	order, err := os.orderRepo.FindByBookID(req.BookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to check if book is ordered."})
		return
	}
	if order != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "This book has already been purchased."})
		return
	}

	newOrder := model.Order{
		UserID:       req.UserID,
		BookID:       req.BookID,
		PurchaseDate: time.Now(),
	}

	err = os.orderRepo.Store(&newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order."})
		return
	}

	book.IsSold = true
	err = os.bookRepo.Update(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update book status."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book purchased successfully.", "order": newOrder})
}

func (os OrderService) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Order ID is invalid."})
		return
	}

	var order model.Order
	err = os.orderRepo.FindByID(id, &order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Order not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (os OrderService) GetUserOrders(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is invalid."})
		return
	}

	orders, err := os.orderRepo.FindAllByUserID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch orders."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
