package db_test

import (
	"github.com/gin-gonic/gin"
	"simple-crud-api/test_db"
	"testing"
)

func TestCreatUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db.DatabaseRefresh()
}
