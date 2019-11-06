package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"offline.com/common"
	"offline.com/service"
)

// ======SERIALIZERS======
type ResourceTypeSerializer struct {
	C *gin.Context
	service.ResourceType
}

type ResourceTypeResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ResourceTypesSerializer struct {
	C             *gin.Context
	ResourceTypes []service.ResourceType
}

func (s *ResourceTypeSerializer) Response() ResourceTypeResponse {
	response := ResourceTypeResponse{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
	}
	return response
}

func (s *ResourceTypesSerializer) Response() []ResourceTypeResponse {
	response := []ResourceTypeResponse{}
	for _, resourceType := range s.ResourceTypes {
		serializer := ResourceTypeSerializer{s.C, resourceType}
		response = append(response, serializer.Response())
	}
	return response
}

func ResourceTypesAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", ResourceTypesList)
}

func ResourceTypesList(c *gin.Context) {
	resourceTypeModels, err := service.GetAllResourceTypes()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("resourceTypes", errors.New("Invalid param")))
		return
	}
	serializer := ResourceTypesSerializer{c, resourceTypeModels}
	c.JSON(http.StatusOK, serializer.Response())
}
