package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/ports"
	"github.com/wansanjou/backend-exercise-user-api/middleware"
)

type userhdl struct {
	usersvc ports.UserService
}

func NewUserHandler(usersvc ports.UserService) *userhdl {
	return &userhdl{
		usersvc: usersvc,
	}
}

func (h *userhdl) UserRoutes(rg *gin.RouterGroup) {
	publicUsers := rg.Group("/users")
	publicUsers.POST("/", h.CreateUser)

	//Authenticated routes
	protectedUsers := rg.Group("/users")
	protectedUsers.Use(middleware.AuthenMiddleware())
	protectedUsers.GET("/", h.GetUsers)
	protectedUsers.GET("/:id", h.GetUserByID)
	protectedUsers.POST("/transfer", h.TransferUser)
}

func (h *userhdl) CreateUser(c *gin.Context) {
	var req domains.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usersvc.CreateUser(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":      user.Name,
		"Email":     user.Email,
		"CreatedAt": user.CreatedAt,
	})
}

func (h *userhdl) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	user, err := h.usersvc.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Name":      user.Name,
		"Email":     user.Email,
		"CreatedAt": user.CreatedAt,
	})
}

func (h *userhdl) GetUsers(c *gin.Context) {
	var req domains.FindAllUsers

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.usersvc.GetUsers(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res []gin.H
	for _, user := range users {
		res = append(res, gin.H{
			"Name":      user.Name,
			"Email":     user.Email,
			"CreatedAt": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (h *userhdl) TransferUser(c *gin.Context) {
	var req domains.TransferRequest
	// เปลี่ยนจาก ShouldBindQuery เป็น ShouldBindJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// เพิ่ม logging เพื่อ debug
	log.Printf("Transfer request: %+v", req)

	err := h.usersvc.TransferBalance(c, req.FromUserID, req.ToUserID, req.Amount)
	if err != nil {
		// เพิ่ม logging error เพื่อ debug
		log.Printf("Transfer error: %v", err)

		// แยก error types เพื่อ return status code ที่เหมาะสม
		if strings.Contains(err.Error(), "insufficient balance") ||
			strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer completed successfully",
		"from":    req.FromUserID,
		"to":      req.ToUserID,
		"amount":  req.Amount,
	})
}
