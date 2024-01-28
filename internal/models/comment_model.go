package models

import (
	"os/user"
	"time"
)

type Comment struct {
	Id        uint   `gorm:"primaryKey"`
	PostId    uint   `gorm:"foreignkey:PostId" json:"postId" binding:"required, gt=0"`
	UserId    uint   `gorm:"foreignkey:UserId"`
	Body      string `gorm:"type:text"`
	CreatedAt time.Time
	User      user.User
}
