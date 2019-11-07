package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"offline.com/common"
	"offline.com/middlewares"
	"offline.com/service"
)

func UsersRegister(router *gin.RouterGroup) {
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
}

func UserRegister(router *gin.RouterGroup) {
	router.GET("/", UserRetrieve)
	router.PUT("/", UserUpdate)
}

func ProfileRegister(router *gin.RouterGroup) {
	router.GET("/me", ProfileRetrieve)
}

func ProfileRetrieve(c *gin.Context) {
	myUserModel := c.MustGet("my_user_model").(service.User)
	userModel, err := service.RetrieveProfile(myUserModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	profileSerializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, profileSerializer.Response())
}

func UsersRegistration(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	if err := service.NewUser(&userModelValidator.userModel, &userModelValidator.cityModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	c.Set("my_user_model", userModelValidator.userModel)
	serializer := UserSerializer{c}
	c.JSON(http.StatusCreated, gin.H{"user": serializer.Response()})
}

func UsersLogin(c *gin.Context) {
	loginValidator := NewLoginValidator()
	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	userModel, err := service.FindOneUser(&service.User{Email: loginValidator.userModel.Email})

	if err != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}

	if userModel.CheckPassword(loginValidator.User.Password) != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}
	middlewares.UpdateContextUserModel(c, userModel.ID)
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func UserRetrieve(c *gin.Context) {
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func UserUpdate(c *gin.Context) {
	myUserModel := c.MustGet("my_user_model").(service.User)
	userModelValidator := NewUserModelValidatorFillWith(myUserModel)
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	userModelValidator.userModel.ID = myUserModel.ID
	if err := myUserModel.Update(userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	middlewares.UpdateContextUserModel(c, myUserModel.ID)
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

// =======SERIALIZERS========

type ProfileSerializer struct {
	C *gin.Context
	service.User
}

// ProfileResponse Declare your response schema here
type ProfileResponse struct {
	ID       uint         `json:"id"`
	Username string       `json:"username"`
	Image    *string      `json:"image"`
	Email    string       `json:"email"`
	City     CityResponse `json:"city"`
}

// Put your response logic including wrap the userModel here.
func (self *ProfileSerializer) Response() ProfileResponse {
	citySerializer := CitySerializer{self.C, self.City}
	fmt.Print()
	if &self.User == nil {
		return ProfileResponse{}
	}
	profile := ProfileResponse{
		ID:       self.ID,
		Username: self.Username,
		Email:    self.Email,
		Image:    self.Image,
		City:     citySerializer.Response(),
	}
	return profile
}

type UserSerializer struct {
	c *gin.Context
}

type UserResponse struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Bio      string  `json:"bio"`
	Image    *string `json:"image"`
	Token    string  `json:"token"`
}

func (self *UserSerializer) Response() UserResponse {
	myUserModel := self.c.MustGet("my_user_model").(service.User)
	user := UserResponse{
		Username: myUserModel.Username,
		Email:    myUserModel.Email,
		Bio:      myUserModel.Bio,
		Image:    myUserModel.Image,
		Token:    common.GenToken(myUserModel.ID),
	}
	return user
}

// ========VALIDATORS=========
// *ModelValidator containing two parts:
// - Validator: write the form/json checking rule according to the doc https://github.com/go-playground/validator
// - DataModel: fill with data from Validator after invoking common.Bind(c, self)
// Then, you can just call model.save() after the data is ready in DataModel.
type UserModelValidator struct {
	User struct {
		Username string `form:"username" json:"username" binding:"exists,alphanum,min=4,max=255"`
		Email    string `form:"email" json:"email" binding:"exists,email"`
		Password string `form:"password" json:"password" binding:"exists,min=8,max=255"`
		Bio      string `form:"bio" json:"bio" binding:"max=1024"`
		Image    string `form:"image" json:"image" binding:"omitempty,url"`
		City     string `form:"city" json:"city" binding:"omitempty"`
	} `json:"user"`
	userModel service.User `json:"-"`
	cityModel service.City `json:"-"`
}

// There are some difference when you create or update a model, you need to fill the DataModel before
// update so that you can use your origin data to cheat the validator.
// BTW, you can put your general binding logic here such as setting password.
func (self *UserModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.userModel.Username = self.User.Username
	self.userModel.Email = self.User.Email
	self.userModel.Bio = self.User.Bio

	if self.User.Password != common.NBRandomPassword {
		self.userModel.SetPassword(self.User.Password)
	}
	if self.User.Image != "" {
		self.userModel.Image = &self.User.Image
	}
	self.cityModel.Name = self.User.City
	return nil
}

// You can put the default value of a Validator here
func NewUserModelValidator() UserModelValidator {
	userModelValidator := UserModelValidator{}
	//userModelValidator.User.Email ="w@g.cn"
	return userModelValidator
}

func NewUserModelValidatorFillWith(userModel service.User) UserModelValidator {
	userModelValidator := NewUserModelValidator()
	userModelValidator.User.Username = userModel.Username
	userModelValidator.User.Email = userModel.Email
	userModelValidator.User.Bio = userModel.Bio
	userModelValidator.User.Password = common.NBRandomPassword

	if userModel.Image != nil {
		userModelValidator.User.Image = *userModel.Image
	}
	return userModelValidator
}

type LoginValidator struct {
	User struct {
		Email    string `form:"email" json:"email" binding:"exists,email"`
		Password string `form:"password"json:"password" binding:"exists,min=8,max=255"`
	} `json:"user"`
	userModel service.User `json:"-"`
}

func (self *LoginValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}

	self.userModel.Email = self.User.Email
	return nil
}

// You can put the default value of a Validator here
func NewLoginValidator() LoginValidator {
	loginValidator := LoginValidator{}
	return loginValidator
}
