package main

import (
	"log"

	"github.com/dan1M/gorm-user-auth-tuto/config"
	"github.com/dan1M/gorm-user-auth-tuto/handler"
	"github.com/dan1M/gorm-user-auth-tuto/model"
	"github.com/dan1M/gorm-user-auth-tuto/service"
	"github.com/gin-gonic/gin"
)

func main() {
	conf := config.InitConfig()
	db, err := config.InitDB(conf)
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&model.User{})

	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	r := gin.Default()

	userApi := r.Group("/api/v1/user")
	userApi.GET("/:id", userHandler.GetUser)
	userApi.GET("/", userHandler.GetUsers)
	userApi.POST("/", userHandler.CreateUser)
	userApi.PUT("/:id", userHandler.UpdateUser)
	userApi.DELETE("/:id", userHandler.DeleteUser)

	r.Run()

}
