package main

import (
	"log"

	"github.com/dan1M/gorm-user-auth-tuto/config"
	"github.com/dan1M/gorm-user-auth-tuto/model"
)

func main() {
	conf := config.InitConfig()
	db, err := config.InitDB(conf)
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&model.User{})

}
