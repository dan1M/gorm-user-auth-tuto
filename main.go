package main

import (
	"fmt"
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
	rtService := service.NewRtService(db)
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(rtService, userService, conf)

	r := gin.Default()

	userApi := r.Group("/api/v1/user")
	userApi.GET("/:id", userHandler.GetUser)
	userApi.GET("/", userHandler.GetUsers)
	userApi.POST("/", userHandler.CreateUser)
	userApi.PUT("/:id", userHandler.UpdateUser)
	userApi.DELETE("/:id", userHandler.DeleteUser)

	authApi := r.Group("/api/v1/auth")
	authApi.POST("/login", authHandler.Login)

	userApi.GET("/test", authHandler.AuthMiddleware(), func(c *gin.Context) {
		fmt.Println(c.Get("user"))
		c.JSON(200, gin.H{"message": "success"})
	})
	r.Run()

}
