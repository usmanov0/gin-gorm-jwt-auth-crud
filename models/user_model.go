package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Email    string `gorm:"column:email;type:varchar(255);unique;not null" json:"email"`
	Password string `gorm:"column:password;type:varchar(255);not null" json:"-"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Email string `json:"email" binding:"required,email"`
}

type SignUpRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type GetUserResponse struct {
	Users []User `json:"users"`
}

type UpdateRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateResponse struct {
	User User `json:"user"`
}
