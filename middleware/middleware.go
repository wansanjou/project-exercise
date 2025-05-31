package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"github.com/wansanjou/backend-exercise-user-api/internal/core/domains"
)

func AuthenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtKey := []byte(viper.GetString("jwt.secretKey"))
		auth := c.Request.Header.Get("Authorization")
		if len(auth) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "token not found",
			})
			c.Abort()
			return
		}

		authPrefix := "Bearer"
		i := len(authPrefix)
		if len(auth) <= i || string(auth[:i]) != authPrefix {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid found",
			})
			c.Abort()
			return
		}
		claims := &domains.JWTClaims{}

		token, err := jwt.ParseWithClaims(string(auth[i+1:]), claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("user", claims)
		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		duration := time.Since(start)

		log.Printf("[HTTP] %s %s %v", method, path, duration)
	}
}
