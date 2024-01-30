package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name  string `gorm:"column:name;unique;not null" json:"name"`
	Slug  string `gorm:"column:slug;unique;not null" json:"slug"`
	Posts []Post `gorm:"foreignKey:CategoryId" json:"posts"`
}
