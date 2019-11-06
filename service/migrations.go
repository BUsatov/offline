package service

import "github.com/jinzhu/gorm"

func Migrate(db *gorm.DB) {

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&ResourceType{})
	db.AutoMigrate(&Event{})
	db.AutoMigrate(&Resource{})
	db.AutoMigrate(&City{})
}
