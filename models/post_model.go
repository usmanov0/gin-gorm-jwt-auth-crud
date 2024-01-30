package models

import (
	"gorm.io/gorm"
	"os/user"
)

type Post struct {
	gorm.Model
	Title      string    `gorm:"column:title;not null" json:"title"`
	Body       string    `gorm:"column:body;not null" json:"body"`
	UserId     uint      `gorm:"column:user_id;not null" json:"user_id"`
	CategoryId uint      `gorm:"column:category_id;not null" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryId" json:"category"`
	User       user.User `gorm:"foreignKey:UserId" json:"user"`
	Comments   []Comment `gorm:"foreignKey:PostId" json:"comments"`
}
