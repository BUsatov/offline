package service

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"offline.com/common"
)

// Models should only be concerned with database schema, more strict checking should be put in validator.
//
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
//
// HINT: If you want to split null and "", you should use *string instead of string.
type User struct {
	ID           uint    `gorm:"primary_key"`
	Username     string  `gorm:"column:username"`
	Email        string  `gorm:"column:email;unique_index"`
	Bio          string  `gorm:"column:bio;size:1024"`
	Image        *string `gorm:"column:image"`
	PasswordHash string  `gorm:"column:password;not null"`
	City         City
	CityID       uint
	Resources    []Resource
	Events       []Event
}

// What's bcrypt? https://en.wikipedia.org/wiki/Bcrypt
// Golang bcrypt doc: https://godoc.org/golang.org/x/crypto/bcrypt
// You can change the value in bcrypt.DefaultCost to adjust the security index.
// 	err := userModel.setPassword("password0")
func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

// Database will only save the hashed string, you should check it by util function.
// 	if err := serModel.checkPassword("password0"); err != nil { password error }
func (u *User) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

// You could input the conditions and it will return an UserModel in database with error info.
// 	userModel, err := FindOneUser(&UserModel{Username: "username0"})
func FindOneUser(condition interface{}) (User, error) {
	db := common.GetDB()
	var model User
	err := db.Where(condition).First(&model).Error
	return model, err
}

func RetrieveProfile(model User) (User, error) {
	db := common.GetDB()
	err := db.Model(&model).Related(&model.City, "City").Error
	return model, err
}

// You could input an UserModel which will be saved in database returning with error info
// 	if err := SaveOne(&userModel); err != nil { ... }
func NewUser(data *User, city *City) error {
	db := common.GetDB()
	var cityModel City
	db.FirstOrCreate(&cityModel, city)
	fmt.Print(city)
	data.CityID = cityModel.ID
	err := db.Save(data).Error
	return err
}

// You could update properties of an UserModel to database returning with error info.
//  err := db.Model(userModel).Update(UserModel{Username: "wangzitian0"}).Error
func (model *User) Update(data interface{}) error {
	db := common.GetDB()
	err := db.Model(model).Update(data).Error
	return err
}
