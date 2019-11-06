package service

import (
	"github.com/jinzhu/gorm"

	"offline.com/common"
)

type City struct {
	gorm.Model
	Name string `gorm:"column:name"`
}

func NewCity(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}

func GetAllCities() ([]City, error) {
	db := common.GetDB()
	var cities []City
	err := db.Find(&cities).Error
	return cities, err
}
