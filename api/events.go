package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"offline.com/common"
	"offline.com/service"
)

func EventsRegister(router *gin.RouterGroup) {
	router.POST("/", EventCreate)
	router.GET("/", EventList)
	router.GET("/:id", EventGet)
}

func EventCreate(c *gin.Context) {
	eventModelValidator := NewEventModelValidator()

	if err := eventModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	if err := service.NewEvent(&eventModelValidator.eventModel, &eventModelValidator.cityModel, &eventModelValidator.resourceModels); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := EventSerializer{c, eventModelValidator.eventModel}
	c.JSON(http.StatusCreated, gin.H{"event": serializer.Response()})
}

func EventList(c *gin.Context) {
	cityID, _ := strconv.ParseUint(c.Query("cityId"), 10, 64)
	categoryID, _ := strconv.ParseUint(c.Query("categoryId"), 10, 64)
	eventsModels, _, err := service.FindManyEvents(uint(cityID), uint(categoryID))
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("events", errors.New("Invalid param")))
		return
	}
	serializer := EventsSerializer{c, eventsModels}
	c.JSON(http.StatusOK, serializer.Response())
}

func EventGet(c *gin.Context) {
	ID := c.Param("id")
	event, err := service.FindEventById(ID)
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("events", errors.New("Invalid param")))
		return
	}
	serializer := EventSerializer{c, event}
	c.JSON(http.StatusOK, serializer.Response())
}

// =======SERIALIZERS=======
type EventSerializer struct {
	C *gin.Context
	service.Event
}

type ResourceResponse struct {
	ID    uint             `json:"id"`
	Value string           `json:"value"`
	User  *ProfileResponse `json:"assignee"`
}

type EventResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	CreatedAt   string             `json:"createdAt"`
	UpdatedAt   string             `json:"updatedAt"`
	User        ProfileResponse    `json:"owner"`
	Category    CategoryResponse   `json:"category"`
	Resources   []ResourceResponse `json:"resources"`
}

type EventsSerializer struct {
	C      *gin.Context
	Events []service.Event
}

type ResourceSerializer struct {
	C *gin.Context
	service.Resource
}

type ResourcesSerializer struct {
	C         *gin.Context
	Resources []service.Resource
}

func (s *ResourceSerializer) Response() ResourceResponse {
	var profileResponse *ProfileResponse = nil
	if s.User != nil {
		userSerializer := ProfileSerializer{s.C, *s.User}
		resp := userSerializer.Response()
		profileResponse = &resp
	}

	response := ResourceResponse{
		ID:    s.ID,
		Value: s.Value,
		User:  profileResponse,
	}
	return response
}

func (s *ResourcesSerializer) Response() []ResourceResponse {
	response := []ResourceResponse{}
	for _, resource := range s.Resources {
		serializer := ResourceSerializer{s.C, resource}
		response = append(response, serializer.Response())
	}
	return response
}

func (s *EventSerializer) Response() EventResponse {
	userSerializer := ProfileSerializer{s.C, s.User}
	categorySetializer := CategorySerializer{s.C, s.Category}
	resourcesSerializer := ResourcesSerializer{s.C, s.Resources}
	response := EventResponse{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		//UpdatedAt:      s.UpdatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt: s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		User:      userSerializer.Response(),
		Category:  categorySetializer.Response(),
		Resources: resourcesSerializer.Response(),
	}

	return response
}

func (s *EventsSerializer) Response() []EventResponse {
	response := []EventResponse{}
	for _, article := range s.Events {
		serializer := EventSerializer{s.C, article}
		response = append(response, serializer.Response())
	}
	return response
}

// ========VALIDATORS=======
type EventModelValidator struct {
	Name        string `form:"name" json:"name" binding:"exists,required"`
	Description string `form:"description" json:"description" binding:"max=2048"`
	CategoryID  uint   `form:"categoryId" json:"categoryId" binding:"required"`
	City        string `form:"city" json:"city" binding:"exists,required"`
	Resources   []struct {
		ResourceTypeID uint   `form:"resourceTypeId" json:"resourceTypeId" binding:"required"`
		Value          string `form:"value" json:"value" binding:"required"`
	} `json:"resources"`
	eventModel     service.Event      `json:"-"`
	cityModel      service.City       `json:"-"`
	resourceModels []service.Resource `json:"-"`
}

func NewEventModelValidator() EventModelValidator {
	return EventModelValidator{}
}

func (s *EventModelValidator) Bind(c *gin.Context) error {
	myUserModel := c.MustGet("my_user_model").(service.User)
	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	s.eventModel.Name = s.Name
	s.eventModel.Description = s.Description
	s.eventModel.CategoryID = s.CategoryID
	s.eventModel.User = myUserModel
	s.cityModel.Name = s.City
	for _, res := range s.Resources {
		var resource service.Resource
		resource.Value = res.Value
		resource.ResourceTypeID = res.ResourceTypeID
		s.resourceModels = append(s.resourceModels, resource)
	}
	return nil
}
