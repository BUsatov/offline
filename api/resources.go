package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"offline.com/common"
	"offline.com/service"
)

func ResourcesRegister(router *gin.RouterGroup) {
	router.PATCH("/:id/assign", AssignResource)
}

func AssignResource(c *gin.Context) {
	ID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	myUserModel := c.MustGet("my_user_model").(service.User)
	resourceModel, err := service.FindResourceById(uint(ID))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("resource", err))
		return
	}

	if err := resourceModel.Assign(myUserModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("assign", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})

	// if err := eventModelValidator.Bind(c); err != nil {
	// 	fmt.Println(err)
	// 	c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
	// 	return
	// }

	// if err := service.NewEvent(&eventModelValidator.eventModel, &eventModelValidator.cityModel, &eventModelValidator.resourceModels); err != nil {
	// 	c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
	// 	return
	// }
	// serializer := EventSerializer{c, eventModelValidator.eventModel}
	// c.JSON(http.StatusCreated, gin.H{"event": serializer.Response()})
}
