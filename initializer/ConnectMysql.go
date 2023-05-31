package initializer

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"fmt"
	"os"
)

func ConnectDataBase() (*gorm.DB ,error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DBUSER"), os.Getenv("DBPASS"), "127.0.0.1", "3306", "BadrCham")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to the database")

	SyncDatabase(db)
	return db, nil;
}