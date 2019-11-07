package service

import (
	"github.com/jinzhu/gorm"

	"offline.com/common"
)

type Event struct {
	gorm.Model
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	User        User
	UserID      uint
	City        City
	CityID      uint
	Category    Category
	CategoryID  uint
	Resources   []Resource
}

func NewEvent(data *Event, cityData interface{}, resourceModels *[]Resource) error {
	db := common.GetDB()
	var cityModel City
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
	tx.Where(Event{CityID: cityID, CategoryID: categoryID}).Preload("Resources.User").Preload("Resources").Preload("Category").Preload("User").Preload("City").Find(&models)

	err := tx.Commit().Error
	return models, count, err
}

func FindEventById(ID string) (Event, error) {
	db := common.GetDB()
	var event Event
	err := db.Where(ID).Preload("User").Preload("Resources").Preload("Resources.User").Preload("Resources.User.City").Preload("Category").Preload("City").Find(&event).Error
	return event, err
}
