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

	c.JSON(http.StatusOK, "ok")
}
