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
	"simple-crud-api/pkg/helper"
	"simple-crud-api/pkg/pagination"
	"simple-crud-api/pkg/util"
	"simple-crud-api/storage/initializers"
	"time"
)

// @Summary Sign up a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User details for sign up"
// @Success 200 {object} gin.H {"message":"Successfully signed up"}
// @Failure 400 {object} gin.H {"validation":{}} "Bad request"
// @Failure 500 {object} gin.H {"error": "Internal Server Error"}
// @Router /api/sign-up [post]
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

// @Summary Sign in a user
// @Description Log in an existing user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.SignInRequest true "User credentials for sign in"
// @Success 200 {object} gin.H {"message":"Successfully logged in"}
// @Failure 400 {object} gin.H {"error": "Invalid email or password"}
// @Failure 500 {object} gin.H {"error": "Internal Server Error"}
// @Router /api/log-in [post]
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

// @Summary Log out the authenticated user
// @Description Log out the currently authenticated user
// @Tags Auth
// @Produce json
// @Success 200 {object} gin.H {"message": "Log out successfully"}
// @Router /api/log-out [post]
func LogOut(c *gin.Context) {
	c.SetCookie("Authorization", "", 0, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "log out successfully",
	})
}

// @Summary Get a list of users
// @Description Retrieve a paginated list of users
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param perPage query int false "Number of users per page"
// @Security Bearer
// @Success 200 {object} gin.H {"result":"models.GetUserResponse"}
// @Failure 401 {object} gin.H {"error": "Unauthorized"}
// @Failure 500 {object} gin.H {"error": "Internal Server Error"}
// @Router /api/users [get]
func GetUsers(c *gin.Context) {
	_, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input pagination.PaginationInput
	if err := c.ShouldBindQuery(&input); err != nil {
		util.HandleValidationErrors(c, err)
		return
	}

	var users []models.User

	res, err := pagination.Paginate(initializers.DB, input.Page, input.Limit, nil, &users)
	if err != nil {
		errors.InternalServerError(c)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": res,
	})
}

// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateRequest true "Updated user details"
// @Security Bearer
// @Success 200 {object} gin.H {"result":"models.UpdateResponse"}
// @Failure 401 {object} gin.H {"error": "Unauthorized"}
// @Failure 403 {object} gin.H {"error": "Forbidden: You are not allowed to update this profile"}
// @Failure 422 {object} gin.H {"validation":{}}
// @Failure 500 {object} gin.H {"error": "Internal Server Error"}
// @Router /api/users/update/{id} [put]
func UpdateUser(c *gin.Context) {
	authUser, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

	if userModel.ID != authUser.Id {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: You are not allowed to update this profile",
		})
		return
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

// @Summary Delete user
// @Description Delete the authenticated user's account
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} gin.H{"message": "User successfully deleted"}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/users/delete/{id} [delete]
func DeleteUser(c *gin.Context) {
	authUser, err := helper.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	id := c.Param("id")

	var user models.User

	result := initializers.DB.First(&user, id)
	if err := result.Error; err != nil {
		errors.RecordNotFound(c, err)
	}

	if user.ID != authUser.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: You are not allowed to delete this profile"})
	}
	initializers.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully deleted",
	})
}
