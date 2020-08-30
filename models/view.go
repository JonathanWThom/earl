package models

import "time"

type View struct {
	ID         uint       `gorm:"primary_key" json:"-"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `sql:"index" json:"-"`
	LinkID     uint       `json:"-"`
	RemoteAddr string     `json:"remoteAddr"`
	UserAgent  string     `json:"userAgent"`
	Referer    string     `json:"referer"`
	Location   Location   `json:"location"`
}
