package middleware

import (
	"fmt"
	"golang-auth/initializers"
	"golang-auth/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid or not provided Token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(401, gin.H{
				"error": "token expired",
			})
			return
		}

		var user models.User
		initializers.DB.First(&user, claims["id"])

		if user.ID == 0 {
			c.JSON(401, gin.H{
				"error": "Invalid token",
			})
			return
		}

		data := map[string]interface{}{
			"ID":    user.ID,
			"Email": user.Email,
		}
		c.Set("user", data)

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
