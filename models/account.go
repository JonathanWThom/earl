package models

import (
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model `json:"-"`
	Token      string `gorm:"unique;not null" json:"token"`
	Links      []Link `json:"-"`
}
