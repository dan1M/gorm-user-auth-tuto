package model

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	User   User
	UserID int    `json:"userId" gorm:"<-:create"`
	Ip     string `json:"ip" gorm:"<-:create"`
	Hash   string `json:"hash" gorm:"<-:create unique index"`
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	rt.CreatedAt = time.Now()
	rt.UpdatedAt = time.Now()

	return
}
