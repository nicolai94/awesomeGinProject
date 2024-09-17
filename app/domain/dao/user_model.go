package dao

import (
	"awesomeProject/app/domain/enums"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string           `gorm:"type:uuid;primary_key" json:"id"`
	Name     string           `gorm:"column:name" json:"name"`
	Email    string           `gorm:"column:email" json:"email"`
	Password string           `gorm:"column:password;size:60" json:"password"`
	Status   enums.UserStatus `gorm:"column:status" json:"status"`
	BaseModel
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}

type UserResponse struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Email  string           `json:"email"`
	Status enums.UserStatus `gorm:"column:status" json:"status"`
}
