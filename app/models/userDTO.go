package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"type:char(36);primaryKey"`
	Username string `json:"username" form:"Username"`
	Email    string `json:"email" form:"Email"`
	Password string `json:"password" form:"Password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
