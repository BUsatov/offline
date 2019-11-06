package service

import (
	"errors"

	"github.com/jinzhu/gorm"
	"offline.com/common"
)

type Resource struct {
	gorm.Model
	Value          string `gorm:"column:value"`
	User           User
	UserID         uint
	Event          Event
	EventID        uint
	ResourceType   ResourceType
	ResourceTypeID uint
}

func (u *Resource) Assign(user User) error {
	db := common.GetDB()
	if u.UserID != 0 {
		return errors.New("already assigned")
	}
	if u.UserID == user.ID {
		return errors.New("this user is already assigned")
	}
	if u.Event.UserID == user.ID {
		return errors.New("event creator can't be assigned")
	}
	err := db.Model(u).Update(Resource{UserID: user.ID}).Error
	if err != nil {
		return errors.New("can't assign user")
	}

	return nil
}

func FindResourceById(ID uint) (Resource, error) {
	db := common.GetDB()
	var resource Resource
	err := db.Where(ID).Preload("Event").Find(&resource).Error
	return resource, err
}
