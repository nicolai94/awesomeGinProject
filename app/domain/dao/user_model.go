package dao

import "awesomeProject/app/domain/enums"

type User struct {
	ID       int              `gorm:"column:id; primary_key; not null" json:"id"`
	Name     string           `gorm:"column:name" json:"name"`
	Email    string           `gorm:"column:email" json:"email"`
	Password string           `gorm:"column:password" json:"-"`
	Status   enums.UserStatus `gorm:"column:status" json:"status"`
	BaseModel
}
