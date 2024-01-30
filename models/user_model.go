package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"column:name;not null" json:"name"`
	Email    string `gorm:"column:email;unique;not null" json:"email"`
	Password string `gorm:"column:password;not null" json:"-"`
}
