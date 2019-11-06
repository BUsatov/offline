package service

import (
	"github.com/jinzhu/gorm"
	"offline.com/common"
)

type Category struct {
	gorm.Model
	Name string `gorm:"unique_index"`
}

func GetAllCategories() ([]Category, error) {
	db := common.GetDB()
	var categories []Category
	err := db.Find(&categories).Error
	return categories, err
}
