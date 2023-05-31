package initializer

import (
	"API/LoginAPI/usermodel"
	"gorm.io/gorm"
)

func SyncDatabase(db *gorm.DB) {
	db.AutoMigrate(&usermodel.Usermodel{})
	
}