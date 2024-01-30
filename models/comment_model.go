package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Body   string `gorm:"column:body" json:"body"`
	PostId uint   `gorm:"foreignKey:PostId" json:"postId" binding:"required, gt=0"`
	UserId uint   `gorm:"foreignKey:UserId"`
	User   User   `gorm:"foreignKey:UserId" json:"user"`
}
