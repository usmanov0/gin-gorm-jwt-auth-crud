package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"simple-crud-api/storage/initializers"
)

func IsUniqueValue(tableName, fieldName, value string) bool {
	var count int64

	result := initializers.DB.Table(tableName).Where(fieldName+" = ?", value).Count(&count)

	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return false
	}

	return count > 0
}

func IsExistValue(tableName, fieldName string, value interface{}) bool {
	var count int64

	result := initializers.DB.Table(tableName).Where(fieldName+" = ?", value).Count(&count)

	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return false
	}

	return count > 0
}

func HandleValidationErrors(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validation": FormatValidationErrors(errs),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

func FormatValidationErrors(errs validator.ValidationErrors) map[string]string {
	errorMessages := make(map[string]string)

	for _, err := range errs {
		fmt.Println()
		switch err.Tag() {
		case "required":
			errorMessages[err.Field()] = fmt.Sprintf("%s is required", err.Field())
		case "email":
			errorMessages[err.Field()] = fmt.Sprintf("%s must be a valid email address", err.Field())
		case "min":
			errorMessages[err.Field()] = fmt.Sprintf("%s must have at least %s characters", err.Field(), err.Param())
		case "max":
			errorMessages[err.Field()] = fmt.Sprintf("%s must have at most %s characters", err.Field(), err.Param())
		case "gt":
			errorMessages[err.Field()] = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
		case "gte":
			errorMessages[err.Field()] = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
		default:
			errorMessages[err.Field()] = fmt.Sprintf("Validation validations on field %s", err.Field())
		}
	}

	return errorMessages
}
