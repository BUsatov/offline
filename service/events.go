package service

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"offline.com/common"
)

type Event struct {
	gorm.Model
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Owner       User
	OwnerID     uint
	City        City
	CityID      uint
	Category    Category
	CategoryID  uint
	Resources   []Resource
}

func NewEvent(data *Event, cityData interface{}, resourceModels *[]Resource) error {
	db := common.GetDB()
	var cityModel City
	fmt.Println(resourceModels)
	db.FirstOrCreate(&cityModel, cityData)
	data.CityID = cityModel.ID
	err := db.Create(&data).Error
	for _, res := range *resourceModels {
		res.EventID = data.ID
		db.Create(&res)
	}
	return err
}

func FindManyEvents(cityID uint, categoryID uint) ([]Event, int, error) {
	db := common.GetDB()
	var models []Event
	var count int

	tx := db.Begin()
	tx.Where(Event{CityID: cityID, CategoryID: categoryID}).Preload("Resources").Preload("Category").Preload("Owner").Preload("City").Find(&models)

	err := tx.Commit().Error
	return models, count, err
}

func FindEventById(ID string) (Event, error) {
	db := common.GetDB()
	var event Event
	err := db.Where(ID).Preload("Resources").Preload("Category").Preload("Owner").Preload("City").Find(&event).Error
	return event, err
}
