package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"simple-crud-api/models"
	"simple-crud-api/pkg/errors"
	"simple-crud-api/pkg/pagination"
	"simple-crud-api/pkg/util"
	"simple-crud-api/storage/initializers"
	"strconv"
	"time"
)

func SignUp(c *gin.Context) {
	var user struct {
		Name     string `json:"name" binding:"required,min=2,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validation": util.FormatValidationErrors(errs),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if util.IsUniqueValue("users", "email", user.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validation": map[string]interface{}{
				"Email": "email already exist!",
			},
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	userModel := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashPassword),
	}

	res := initializers.DB.Create(&userModel)

	if res.Error != nil {
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userModel,
	})
}

func SignIn(c *gin.Context) {
	var user struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if c.ShouldBindJSON(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var userModel models.User
	initializers.DB.First(&userModel, "email = ?", user.Email)

	if userModel.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userModel.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenStr, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func LogOut(c *gin.Context) {
	c.SetCookie("Authorization", "", 0, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "log out successfully",
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "5")
	perPage, _ := strconv.Atoi(perPageStr)

	res, err := pagination.Paginate(initializers.DB, page, perPage, nil, &users)
	if err != nil {
		errors.InternalServerError(c)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": res,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user struct {
		Name  string `json:"name" binding:"required,min=2,max=50"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validation": util.FormatValidationErrors(errs),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var userModel models.User
	res := initializers.DB.First(&userModel, id)

	if err := res.Error; err != nil {
		errors.RecordNotFound(c, err)
	}

	if userModel.Email != user.Email && util.IsUniqueValue("users", "email", user.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"Email": "email is already exist",
			},
		})
		return
	}

	updateUser := models.User{
		Name:  user.Name,
		Email: user.Email,
	}

	result := initializers.DB.Model(&user).Updates(&updateUser)

	if result.Error != nil {
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	result := initializers.DB.First(&user, id)
	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
	}

	initializers.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully deleted",
	})
}
