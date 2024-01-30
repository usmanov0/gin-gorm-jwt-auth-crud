package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-crud-api/middleware"
)

func GetAuthUser(c *gin.Context) (*middleware.AuthUser, error) {
	authUser, exists := c.Get("authUser")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get the user",
		})
		return nil, fmt.Sprintf("Failed to get user")
	}

	if user, ok := authUser.(middleware.AuthUser); ok {
		return &user, nil
	}

	return nil, nil
}
