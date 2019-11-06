package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"offline.com/api"
	"offline.com/common"
	"offline.com/middlewares"
	"offline.com/service"
)

func main() {

	db := common.Init()
	service.Migrate(db)
	defer db.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")
	api.UsersRegister(v1.Group("/users"))
	v1.Use(middlewares.AuthMiddleware(false))

	v1.Use(middlewares.AuthMiddleware(true))
	api.UserRegister(v1.Group("/user"))
	api.ProfileRegister(v1.Group("/profiles"))

	api.CategoriesAnonymousRegister(v1.Group("/categories"))
	api.ResourceTypesAnonymousRegister(v1.Group("/resource-types"))

	api.EventsRegister(v1.Group("/events"))

	api.CitiesRegister(v1.Group("/cities"))
	api.ResourcesRegister(v1.Group("/resources"))

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// test 1 to 1
	tx1 := db.Begin()
	userA := service.User{
		Username: "AAAAAAAAAAAAAAAA",
		Email:    "aaaa@g.cn",
		Bio:      "hehddeda",
		Image:    nil,
	}
	tx1.Save(&userA)
	// tx1.Commit()
	category1 := service.Category{
		Name: "cooking",
	}
	category2 := service.Category{
		Name: "handcrafting",
	}
	category3 := service.Category{
		Name: "sport",
	}
	category4 := service.Category{
		Name: "wellness",
	}
	category5 := service.Category{
		Name: "gardening",
	}
	category6 := service.Category{
		Name: "fixing",
	}
	category7 := service.Category{
		Name: "arts",
	}
	category8 := service.Category{
		Name: "parenting",
	}
	resourceType := service.ResourceType{
		Name: "skills",
	}
	resourceType2 := service.ResourceType{
		Name: "materials",
	}
	resourceType3 := service.ResourceType{
		Name: "location",
	}
	resourceType4 := service.ResourceType{
		Name: "services",
	}
	resourceType5 := service.ResourceType{
		Name: "other",
	}
	tx1.Save(&category1)
	tx1.Save(&category2)
	tx1.Save(&category3)
	tx1.Save(&category4)
	tx1.Save(&category5)
	tx1.Save(&category6)
	tx1.Save(&category7)
	tx1.Save(&category8)
	tx1.Save(&resourceType)
	tx1.Save(&resourceType2)
	tx1.Save(&resourceType3)
	tx1.Save(&resourceType4)
	tx1.Save(&resourceType5)
	tx1.Commit()
	fmt.Println(userA)

	//db.Save(&ArticleUserModel{
	//    UserModelID:userA.ID,
	//})
	//var userAA ArticleUserModel
	//db.Where(&ArticleUserModel{
	//    UserModelID:userA.ID,
	//}).First(&userAA)
	//fmt.Println(userAA)

	r.Run() // listen and serve on 0.0.0.0:8080
}
