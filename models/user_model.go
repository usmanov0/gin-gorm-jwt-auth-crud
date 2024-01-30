package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Email    string `gorm:"column:email;type:varchar(255);unique;not null" json:"email"`
	Password string `gorm:"column:password;type:varchar(255);not null" json:"-"`
}
