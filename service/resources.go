package service

import (
	"errors"

	"github.com/jinzhu/gorm"
	"offline.com/common"
)

type Resource struct {
	gorm.Model
	Value          string `gorm:"column:value"`
	Assignee       User
	AssigneeID     uint
	Event          Event
	EventID        uint
	ResourceType   ResourceType
	ResourceTypeID uint
}

func (u *Resource) Assign(user User) error {
	db := common.GetDB()
	if u.AssigneeID != 0 {
		return errors.New("password should not be empty!")
	}
	if u.AssigneeID == user.ID {
		return errors.New("this user is already assigned")
	}
	err := db.Model(u).Update(Resource{AssigneeID: user.ID}).Error
	if err != nil {
		return errors.New("can't assign user")
	}

	return nil
}

func FindResourceById(ID uint) (Resource, error) {
	db := common.GetDB()
	var resource Resource
	err := db.Where(ID).Find(&resource).Error
	return resource, err
}
