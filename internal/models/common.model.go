package models

import "time"

type CommonModel struct {
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"-"`
	DeletedAt *time.Time `gorm:"default:null"   json:"-"`
}
