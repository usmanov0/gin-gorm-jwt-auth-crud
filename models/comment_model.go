package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Body   string `gorm:"column:body;type:text" json:"body"`
	PostId uint   `gorm:"foreignKey:PostId;type:integer;not null" json:"post_id" binding:"required, gt=0"`
	UserId uint   `gorm:"foreignKey:UserId;type:integer" json:"user_id"`
	User   User   `gorm:"foreignKey:UserId" json:"user"`
}
