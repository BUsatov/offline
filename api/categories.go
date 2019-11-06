package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"offline.com/common"
	"offline.com/service"
)

func CategoriesAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", CategoriesList)
}

func CategoriesList(c *gin.Context) {
	categoryModels, err := service.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("categories", errors.New("Invalid param")))
		return
	}
	serializer := CategoriesSerializer{c, categoryModels}
	c.JSON(http.StatusOK, serializer.Response())
}

// =====SERIALIZERS======
type CategorySerializer struct {
	C *gin.Context
	service.Category
}

type CategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CategoriesSerializer struct {
	C          *gin.Context
	Categories []service.Category
}

func (s *CategorySerializer) Response() CategoryResponse {
	response := CategoryResponse{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		UpdatedAt: s.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
	}
	return response
}

func (s *CategoriesSerializer) Response() []CategoryResponse {
	response := []CategoryResponse{}
	for _, category := range s.Categories {
		serializer := CategorySerializer{s.C, category}
		response = append(response, serializer.Response())
	}
	return response
}
