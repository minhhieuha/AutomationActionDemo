package http

import (
	"fmt"
	"net/http"
	"strconv"
	"testing-demo/domain"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUsecase domain.OrderUsecase
}

// NewOrderHandler creates new order HTTP handler
func NewOrderHandler(router *gin.Engine, uc domain.OrderUsecase) {
	handler := &OrderHandler{
		orderUsecase: uc,
	}

	router.GET("/ping", handler.Ping)
	router.POST("/orders", handler.CreateOrder)
	router.GET("/orders/:id", handler.GetOrder)
}

// Ping health check endpoint
func (h *OrderHandler) Ping(c *gin.Context) {
	fmt.Println("Auto Merge 2 - Source Tree CD - Jenkins")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		CustomerName string `json:"customer_name" binding:"required"`
		Item         string `json:"item" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	order, err := h.orderUsecase.CreateOrder(req.CustomerName, req.Item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Order created successfully",
		"data":    order,
	})
}

// GetOrder handles GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid order ID",
		})
		return
	}

	order, err := h.orderUsecase.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}
