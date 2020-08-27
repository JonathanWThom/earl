package models

import (
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Token string `gorm:"unique;not null"`
	Links []Link
}
