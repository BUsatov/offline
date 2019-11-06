package service

import (
	"github.com/jinzhu/gorm"
	"offline.com/common"
)

type ResourceType struct {
	gorm.Model
	Name string `gorm:"unique_index"`
}

func GetAllResourceTypes() ([]ResourceType, error) {
	db := common.GetDB()
	var resourceTypes []ResourceType
	err := db.Find(&resourceTypes).Error
	return resourceTypes, err
}
