package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/ports"
)

type authhandler struct {
	authsvc ports.AuthService
	usersvc ports.UserService
}

func NewAuthHandler(authsvc ports.AuthService) *authhandler {
	return &authhandler{
		authsvc: authsvc,
	}
}

func (h *authhandler) AuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", h.LoginHandler)
}

func (h *authhandler) LoginHandler(c *gin.Context) {
	var req domains.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.authsvc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
