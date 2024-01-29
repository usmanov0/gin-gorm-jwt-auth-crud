package models

import (
	"gorm.io/gorm"
	"os/user"
)

type Post struct {
	gorm.Model
	Title      string    `gorm:"not null" json:"title"`
	Body       string    `gorm:"not null" json:"body"`
	UserId     uint      `gorm:"foreignkey:UserId" json:"userId"`
	CategoryId uint      `gorm:"foreignkey:CategoryId" json:"categoryId"`
	Category   Category  `gorm:"foreignkey:CategoryId"`
	User       user.User `gorm:"foreignkey:UserId"`
	Comments   []Comment
}
