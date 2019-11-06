package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"offline.com/common"
	"offline.com/service"
)

// =======SERIALIZERS========
type CityResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CitySerializer struct {
	C *gin.Context
	service.City
}

type CitiesSerializer struct {
	C      *gin.Context
	Cities []service.City
}

func (s *CitySerializer) Response() CityResponse {
	response := CityResponse{
		ID:        s.City.Model.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
	}
	return response
}

func (s *CitiesSerializer) Response() []CityResponse {
	response := []CityResponse{}
	for _, resourceType := range s.Cities {
		serializer := CitySerializer{s.C, resourceType}
		response = append(response, serializer.Response())
	}
	return response
}

// ========VALIDATORS========
type CityModelValidator struct {
	Name        string       `form:"name" json:"name" binding:"exists,required"`
	CityModelIn service.City `json:"-"`
}

func NewCityModelValidator() CityModelValidator {
	return CityModelValidator{}
}

func (s *CityModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	s.CityModelIn.Name = s.Name
	return nil
}

func CitiesRegister(router *gin.RouterGroup) {
	router.GET("/", CitiesList)
}

func CitiesList(c *gin.Context) {
	cityModels, err := service.GetAllCities()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("resourceTypes", errors.New("Invalid param")))
		return
	}
	serializer := CitiesSerializer{c, cityModels}
	c.JSON(http.StatusOK, serializer.Response())
}
