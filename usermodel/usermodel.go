package usermodel

import (
	"gorm.io/gorm"
)

type Usermodel struct {
	gorm.Model
	Login		string `gorm: "unique"`
	Password	string
	Admin   bool
}

func (Usermodel) TableName() string {
	return "Users"
}

//   equals
//   type User struct {
// 	ID        uint           `gorm:"primaryKey"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
// 	login string
//   }